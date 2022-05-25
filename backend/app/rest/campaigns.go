package api

import (
	"amifactory.team/sequel/coton-app-backend/app/logger"
	"amifactory.team/sequel/coton-app-backend/app/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var placeholderRegex = regexp.MustCompile(`{{([a-zA-Z_]+)}}`)

type userCampaigns struct {
	contactsStore model.ContactsStore
	pricesStore   model.PricesStore
	phoneService  model.PhoneService
	accountStore  model.AccountStore
	campaignStore model.CampaignStore

	secureCookiesDisabled bool
}

func (c *userCampaigns) add(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, _ := model.UserFromContext(ctx, model.KCtxKeyUserMe)
	account, _ := model.AccountFromContext(ctx)
	member, err := c.accountStore.FindUserAccountMember(ctx, account.ID, user.ID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			httpResponseError(w, http.StatusInternalServerError, "internal")
		} else {
			httpResponseError(w, http.StatusInternalServerError, "internal")
		}

		return
	}

	form, err := newCampaignForm(r)
	if err != nil {
		if errors.Is(err, FormError) {
			httpResponseError(w, http.StatusBadRequest, err.Error())
		} else {
			c.getLog(ctx).Errorf("Fail to decode body - %v", err)
			httpResponseError(w, http.StatusBadRequest, "body_json__invalid")
		}
		return
	}

	errs := form.Validate()
	if len(errs) > 0 {
		if len(form.PhoneErrors) > 0 {
			resp := newCampaignError(errs, form.PhoneErrors)
			httpJsonResponseCode(w, resp, http.StatusBadRequest)
		} else {
			httpResponseErrors(w, http.StatusBadRequest, errs)
		}
		return
	}

	campaign := form.Campaign(account.ID, member.ID)
	campaign.Currency = account.Balance.Currency

	// 1. For 'manual' target option we need to group phones by mcc/mnc codes,
	//	  get prices for each mcc/mnc pair and calculate total amount.
	//    We assume that variables for 'manual' target are not be used
	// 2. For 'group' target option we need to fetch all contacts for selected
	//	  group, compute sms parts for each message (with variables applied),
	//    group sms parts by mcc/mnc codes, get prices for each mcc/mnc pair and
	//    calculate total amount.
	if campaign.Target == model.Manual {
		phoneGroups, err := campaign.GroupPhones()
		if err != nil {
			httpResponseError(w, http.StatusInternalServerError, "internal")
			return
		}

		smsParts := 1

		routes := make([]model.CampaignRoute, 0, len(phoneGroups))
		for mccMnc, count := range phoneGroups {
			price, err := c.pricesStore.FetchPriceForAccount(ctx, account.ID, mccMnc)
			if err != nil {
				if errors.Is(err, model.ErrNotFound) {
					// no problem here, we will mark this route as 'unavailable'
				} else {
					httpResponseError(w, http.StatusInternalServerError, "internal")
					return
				}
			}

			refBookRec, err := c.phoneService.RefBookRecordForMccMnc(ctx, mccMnc)
			if err != nil {
				httpResponseError(w, http.StatusInternalServerError, "internal")
				return
			}

			var route model.CampaignRoute
			smsTotal := smsParts * count
			if price != nil {
				total := price.Price.Multi(smsTotal)
				route = model.CampaignRoute{
					ID:         uuid.NewV4().String(),
					CampaignID: campaign.ID,
					Country:    refBookRec.Country,
					MccMnc:     mccMnc,
					SmsCount:   smsTotal,
					Currency:   price.Currency,
					Price:      price.Price,
					Total:      total,
					Available:  price.Enabled,
				}

			} else {
				route = model.CampaignRoute{
					ID:         uuid.NewV4().String(),
					CampaignID: campaign.ID,
					Country:    refBookRec.Country,
					MccMnc:     mccMnc,
					SmsCount:   smsTotal,
					Price:      model.Zero,
					Total:      model.Zero,
				}
			}
			routes = append(routes, route)
		}

		campaign.SetRoutes(routes)
	}

	if campaign.Target == model.Group {
		groupId := form.GroupID
		group, err := c.contactsStore.FindGroup(ctx, account.ID, groupId)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				httpResponseError(w, http.StatusBadRequest, "group__invalid")
			} else {
				c.getLog(ctx).WithError(err).Error("Fail to fetch group")
				httpResponseError(w, http.StatusInternalServerError, "internal")
			}
			return
		}

		submatches := placeholderRegex.FindAllStringSubmatch(campaign.MessageText, -1)
		placeholdersSet := make(map[string]bool, 0)
		for _, v := range submatches {
			placeholdersSet[v[1]] = true
		}
		placeholders := make([]string, 0)
		for p, _ := range placeholdersSet {
			placeholders = append(placeholders, p)
		}
		vars := group.FindVariablesByName(placeholders)

		contacts, err := c.contactsStore.ListContactsAll(ctx, account.ID, *group)
		if err != nil {
			c.getLog(ctx).WithError(err).Error("Fail to fetch group contacts")
			httpResponseError(w, http.StatusInternalServerError, "internal")
			return
		}

		groupedContacts := make(map[model.MccMnc]int)
		for _, contact := range contacts {
			varValues := make(map[string]string, 0)
			for _, v := range vars {
				value, found := contact.VarValues[v.Name]
				if found {
					varValues[v.Name] = fmt.Sprintf("%v", value)
				}
			}

			message := campaign.MessageText
			for placeholder, value := range varValues {
				message = strings.ReplaceAll(message, fmt.Sprintf("{{%s}}", placeholder), value)
			}

			// TODO compute sms count
			count := groupedContacts[contact.MccMnc]
			groupedContacts[contact.MccMnc] = count + 1
		}

		routes := make([]model.CampaignRoute, 0, len(groupedContacts))
		for mccMnc, count := range groupedContacts {
			price, err := c.pricesStore.FetchPriceForAccount(ctx, account.ID, mccMnc)
			if err != nil {
				if errors.Is(err, model.ErrNotFound) {
					// no problem here, we will mark this route as 'unavailable'
				} else {
					httpResponseError(w, http.StatusInternalServerError, "internal")
					return
				}
			}

			refBookRec, err := c.phoneService.RefBookRecordForMccMnc(ctx, mccMnc)
			if err != nil {
				c.getLog(ctx).WithError(err).Errorf("Fail to find refbook record for %v", mccMnc)
				httpResponseError(w, http.StatusInternalServerError, "internal")
				return
			}

			var route model.CampaignRoute
			if price != nil {
				total := price.Price.Multi(count)
				route = model.CampaignRoute{
					ID:         uuid.NewV4().String(),
					CampaignID: campaign.ID,
					Country:    refBookRec.Country,
					MccMnc:     mccMnc,
					SmsCount:   count,
					Currency:   price.Currency,
					Price:      price.Price,
					Total:      total,
					Available:  price.Enabled,
				}

			} else {
				route = model.CampaignRoute{
					ID:         uuid.NewV4().String(),
					CampaignID: campaign.ID,
					Country:    refBookRec.Country,
					MccMnc:     mccMnc,
					SmsCount:   count,
					Price:      model.Zero,
					Total:      model.Zero,
				}
			}
			routes = append(routes, route)
		}

		campaign.SetRoutes(routes)
	}

	if !form.DryRun {
		err = c.campaignStore.Add(ctx, campaign)
		if err != nil {
			c.getLog(ctx).WithError(err).Errorf("Fail to save campaign")
			httpResponseError(w, http.StatusInternalServerError, "internal")
			return
		}
	}

	httpJsonResponse(w, newCampaignResponse(campaign))
	return
}

func (c *userCampaigns) constraints(w http.ResponseWriter, r *http.Request) {
	latestDatetime := time.Now().AddDate(0, 6, 0).Unix()

	var resp = struct {
		ScheduledCampaignsLeft           int   `json:"scheduled_campaigns_left"`
		ScheduledCampaignsLatestDatetime int64 `json:"scheduled_campaigns_latest_datetime"`
		GroupsMaxNumber                  int   `json:"groups_max_number"`
		PhonesMaxNumber                  int   `json:"phones_max_number"`
	}{
		ScheduledCampaignsLeft:           5,
		ScheduledCampaignsLatestDatetime: latestDatetime * 1000,
		GroupsMaxNumber:                  20,
		PhonesMaxNumber:                  500,
	}

	httpJsonResponse(w, resp)
}

func (c *userCampaigns) getLog(ctx context.Context) *logrus.Entry {
	return logger.GetLogger(ctx).WithFields(logrus.Fields{
		"scope": "user_campaigns",
	})
}

var senderIDRegex = regexp.MustCompile(`^[\dA-Za-z]{1,11}$`)

func isValidSenderID(value string) bool {
	return senderIDRegex.MatchString(value)
}

type campaignForm struct {
	Name        string
	TargetsType model.CampaignTarget

	SenderID    string
	GroupID     string
	Phones      []*model.Phone
	PhoneErrors []string

	MessageText      string
	SendOption       model.SendOption
	SendAfter        *time.Time
	TrackLinkOpening bool

	DryRun bool
}

func newCampaignForm(r *http.Request) (*campaignForm, error) {
	body := struct {
		Name        string `json:"name"`
		TargetsType string `json:"targets_type"`

		SenderID string   `json:"sender_id"`
		GroupID  string   `json:"group"`
		Phones   []string `json:"phones"`

		MessageText      string `json:"message_text"`
		SendOption       string `json:"send_option"`
		SendAfter        *int64 `json:"send_after"`
		TrackLinkOpening bool   `json:"track_link_opening"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	var form campaignForm

	form.Name = ClearString(body.Name)
	if len(form.Name) == 0 {
		form.Name = fmt.Sprintf("Campaign %s", time.Now().Format("2006-01-02 15:04"))
	}
	targetsType, err := model.ParseCampaignTarget(body.TargetsType)
	if err != nil {
		return nil, NewFormError("targets_type__invalid", err)
	}
	form.TargetsType = targetsType
	if form.TargetsType == model.Group {
		form.GroupID = ClearString(body.GroupID)
	}
	if form.TargetsType == model.Manual && len(body.Phones) > 0 {
		form.Phones = make([]*model.Phone, 0)
		form.PhoneErrors = make([]string, 0)
		for idx, phoneStr := range body.Phones {
			phone, err := model.ParsePhone(phoneStr)
			if err != nil {
				if errors.Is(err, model.ErrPhoneInvalidSymbols) {
					form.PhoneErrors = append(form.PhoneErrors, fmt.Sprintf("%d: Phone number '%s' contains invalid symbols", idx+1, phoneStr))
				} else if errors.Is(err, model.ErrPhoneMccMncNotFound) {
					form.PhoneErrors = append(form.PhoneErrors, fmt.Sprintf("%d: Phone number '%s' is not valid - can not find mcc/mnc codes", idx+1, phoneStr))
				} else if errors.Is(err, model.ErrPhoneWrongLen) {
					form.PhoneErrors = append(form.PhoneErrors, fmt.Sprintf("%d: Phone number '%s' lenght is not valid", idx+1, phoneStr))
				} else {
					form.PhoneErrors = append(form.PhoneErrors, fmt.Sprintf("%d: Phone number '%s' in not valid", idx+1, phoneStr))
				}
			} else {
				form.Phones = append(form.Phones, phone)
			}
		}
	}
	form.SenderID = ClearString(body.SenderID)
	form.MessageText = strings.TrimSpace(body.MessageText)
	sendOption, err := model.ParseSendOption(body.SendOption)
	if err != nil {
		return nil, NewFormError("send_option__invalid", err)
	}
	form.SendOption = sendOption
	if body.SendAfter != nil {
		sendAfter := time.Unix(int64(*body.SendAfter/1000), 0)
		form.SendAfter = &sendAfter
	}
	form.TrackLinkOpening = body.TrackLinkOpening

	dryRun, err := strconv.ParseBool(r.URL.Query().Get("dry_run"))
	if err != nil {
		return nil, NewFormError("dry_run__invalid", err)
	}
	form.DryRun = dryRun
	return &form, nil
}

func (f *campaignForm) Validate() []error {
	errs := make([]error, 0)

	if len(f.Name) == 0 {
		errs = append(errs, errors.New("name__required"))
	} else if utf8.RuneCountInString(f.Name) > 150 {
		errs = append(errs, errors.New("name__invalid"))
	}

	if f.TargetsType == model.Group {
		if !IsValidUUID(f.GroupID) {
			errs = append(errs, errors.New("group__invalid"))
		}
	}

	if f.TargetsType == model.Manual {
		if len(f.Phones) == 0 && len(f.PhoneErrors) == 0 {
			errs = append(errs, errors.New("phones__required"))
		} else if len(f.PhoneErrors) > 0 {
			errs = append(errs, errors.New("phones__invalid"))
		}
	}

	if len(f.SenderID) == 0 {
		errs = append(errs, errors.New("sender_id__required"))
	} else if !isValidSenderID(f.SenderID) {
		errs = append(errs, errors.New("sender_id__invalid"))
	}

	if len(f.MessageText) == 0 {
		errs = append(errs, errors.New("message_text__required"))
	} else {
		// TODO message text validation
	}

	if f.SendOption == model.Scheduled {
		if f.SendAfter == nil {
			errs = append(errs, errors.New("send_after__required"))
		} else if f.SendAfter.Before(time.Now()) {
			errs = append(errs, errors.New("send_after__in_past"))
		}
	}

	return errs
}

func (f *campaignForm) Campaign(accountID, memberID string) *model.Campaign {
	c := &model.Campaign{
		ID:               uuid.NewV4().String(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		AccountID:        accountID,
		AccountMemberID:  memberID,
		Name:             f.Name,
		SenderID:         f.SenderID,
		MessageText:      f.MessageText,
		Target:           f.TargetsType,
		SendOption:       f.SendOption,
		SendAfter:        f.SendAfter,
		TrackLinkOpening: f.TrackLinkOpening,
	}
	switch f.TargetsType {
	case model.Manual:
		c.Phones = f.Phones
	case model.Group:
		c.GroupID = f.GroupID
	}
	return c
}

type restCampaignErrorPayload struct {
	restResponseErrorPayload
	PhoneErrors []string `json:"phone_error"`
}

func newCampaignError(errs []error, phoneErrors []string) restCampaignErrorPayload {
	errsStr := make([]string, len(errs))
	for idx, e := range errs {
		errsStr[idx] = e.Error()
	}

	resp := restCampaignErrorPayload{
		restResponseErrorPayload: restResponseErrorPayload{
			Errors: errsStr,
		},
		PhoneErrors: phoneErrors,
	}

	return resp
}

type campaignRouteResponse struct {
	model.CampaignRoute
	FormattedPrice string `json:"price"`
	FormattedTotal string `json:"total"`
}

func newCampaignRouteResponse(route model.CampaignRoute) campaignRouteResponse {
	return campaignRouteResponse{
		CampaignRoute:  route,
		FormattedPrice: route.Price.Format(route.Currency),
		FormattedTotal: route.Total.Format(route.Currency),
	}
}

type campaignResponse struct {
	*model.Campaign
	RoutesResponse []campaignRouteResponse `json:"routes"`
	FormattedTotal string                  `json:"total"`
}

func newCampaignResponse(campaign *model.Campaign) campaignResponse {
	routesResp := make([]campaignRouteResponse, len(campaign.Routes))
	for idx, route := range campaign.Routes {
		routesResp[idx] = newCampaignRouteResponse(route)
	}
	return campaignResponse{
		Campaign:       campaign,
		RoutesResponse: routesResp,
		FormattedTotal: campaign.Total.Format(campaign.Currency),
	}
}
