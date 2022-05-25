package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"amifactory.team/sequel/coton-app-backend/app/logger"
	"amifactory.team/sequel/coton-app-backend/app/model"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

type staffPrices struct {
	pricesStore  model.PricesStore
	accountStore model.AccountStore
	phoneService model.PhoneService
}

func (p *staffPrices) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	form := newPricesPageForm(r)
	if form.HasValidationErrors() {
		httpResponseErrors(w, http.StatusBadRequest, form.ValidationErrors)
		return
	}

	filter := model.PricesFilter{
		SearchPhase: form.SearchPhrase,
		AccountID:   form.AccountID,
	}

	prices, total, err := p.pricesStore.FetchPricesForAccount(ctx, form.Offset(), form.PageSize, filter)
	if err != nil {
		if errors.Is(err, model.ErrPageOutOfBounds) {
			httpResponseError(w, http.StatusBadRequest, "page__out_of_bounds")
		} else {
			p.getLog(ctx).Errorf("Fail to fetch prices list - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
		}

		return
	}

	if total == 0 && form.HasFilter() {
		httpResponseError(w, http.StatusNotFound, "result__not_found")
		return
	}

	isGeneralPrice := form.AccountID == nil

	resp := newStaffPricesResponse(prices, isGeneralPrice)
	httpJsonResponse(w, NewPageResp(total, form.PageSize, resp))
}

func (p *staffPrices) add(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	staff, err := model.StaffFromContext(ctx, model.KCtxKeyStaffMe)
	if err != nil {
		p.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "request__unauthorized")
		return
	}

	// TODO update later
	if staff.Role == nil || staff.Role.Name != "admin" {
		p.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	form := newPriceForm(r)
	if form.HasValidationErrors() {
		httpResponseErrors(w, http.StatusBadRequest, form.ValidationErrors)
		return
	}

	if form.AccountId != nil {
		account, err := p.accountStore.FindAccountByID(ctx, *form.AccountId)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				httpResponseError(w, http.StatusNotFound, "account__not_found")
			} else {
				p.getLog(ctx).WithError(err).Error("Fail to find account")
				httpResponseError(w, http.StatusInternalServerError, "internal")
			}
			return
		}

		if account.Billing.Currency != form.Currency {
			httpResponseError(w, http.StatusBadRequest, "currency__not_match_account")
			return
		}
	}

	price := form.newPrice()
	refBook, err := p.phoneService.RefBookRecordForMccMnc(ctx, price.MccMnc)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			// ok, no problem here
		} else {
			httpResponseError(w, http.StatusInternalServerError, "internal")
			return
		}
	} else {
		price.Country = refBook.Country
		price.Operator = refBook.Operator
	}

	if price.AccountId != nil {
		// check if price item in general price exists
		_, err = p.pricesStore.FetchPriceForAccount(ctx, *price.AccountId, price.MccMnc)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				httpResponseError(w, http.StatusForbidden, "general_price__not_exist")
			} else {
				httpResponseError(w, http.StatusInternalServerError, "internal")
			}
			return
		}
	}

	err = p.pricesStore.AddPrice(ctx, price)
	if err != nil {
		if errors.Is(err, model.ErrDuplicate) {
			httpResponseError(w, http.StatusForbidden, "price__already_exist")
		} else {
			p.getLog(ctx).WithError(err).Error("Fail to save price")
			httpResponseError(w, http.StatusInternalServerError, "internal")
		}
		return
	}

	isGeneralPrice := price.AccountId == nil
	httpJsonResponse(w, newStaffPriceResponse(price, isGeneralPrice))
}

func (p *staffPrices) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	form := newUpdatePriceForm(r)
	if form.HasValidationErrors() {
		httpResponseErrors(w, http.StatusBadRequest, form.ValidationErrors)
		return
	}
	price, _ := model.PriceItemFromContext(ctx)
	form.updatePriceItem(price)
	err := p.pricesStore.UpdatePrice(ctx, *price)
	if err != nil {
		p.getLog(ctx).Errorf("Fail to update price - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	isGeneralPrice := price.AccountId == nil
	httpJsonResponse(w, newStaffPriceResponse(*price, isGeneralPrice))
}

func (p *staffPrices) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	price, _ := model.PriceItemFromContext(ctx)
	err := p.pricesStore.DeletePrice(ctx, *price)
	if err != nil {
		p.getLog(ctx).Errorf("Fail to delete price  - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}
	httpJsonResponse(w, http.StatusNoContent)
}

func (p *staffPrices) populatePriceItem(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		priceID := chi.URLParam(r, "priceID")
		if len(priceID) == 0 {
			httpResponseError(w, http.StatusBadRequest, "price__required")
			return
		}

		price, err := p.pricesStore.FetchPrice(ctx, priceID)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				httpResponseError(w, http.StatusNotFound, "price_item__not_found")
			} else {
				httpResponseError(w, http.StatusInternalServerError, "internal")
			}
			return
		}

		next.ServeHTTP(w, r.WithContext(price.NewContext(ctx)))
	})
}

func (p *staffPrices) getLog(ctx context.Context) *logrus.Entry {
	return logger.GetLogger(ctx).WithFields(logrus.Fields{
		"scope": "staff-prices",
	})
}

/*
 Staff prices forms
*/
type pricesPageForm struct {
	*PageForm
	AccountID *string
}

func newPricesPageForm(r *http.Request) *pricesPageForm {
	pageForm := NewPageForm(r, 1)
	f := pricesPageForm{
		PageForm: pageForm,
	}

	accountIdParam := r.URL.Query().Get("account")

	if len(accountIdParam) > 0 {
		if IsValidUUID(accountIdParam) {
			f.AccountID = &accountIdParam
		} else if accountIdParam == "general" {
			// do nothing, accountID = nil
		} else {
			f.AddValidationError("account__invalid", nil)
		}
	}

	return &f
}

type priceForm struct {
	*BaseForm

	AccountId *string

	Mcc int
	Mnc int

	Price    model.Cost
	Currency model.Currency
}

func newPriceForm(r *http.Request) priceForm {
	body := struct {
		AccountId *string `json:"account"`

		Mcc *int `json:"mcc"`
		Mnc *int `json:"mnc"`

		Price    *int    `json:"price"`
		Currency *string `json:"currency"`
	}{}

	f := priceForm{
		BaseForm: &BaseForm{},
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		f.AddValidationError("body_json__invalid", err)
		return f
	}

	if body.AccountId != nil {
		accountId := ClearString(*body.AccountId)
		if accountId == "general" {
			f.AccountId = nil
		} else if IsValidUUID(accountId) {
			f.AccountId = &accountId
		} else {
			f.AddValidationError("account__invalid", nil)
		}
	}

	if body.Mcc != nil {
		if *body.Mcc >= 0 {
			f.Mcc = *body.Mcc
		} else {
			f.AddValidationError("mcc__invalid", nil)
		}
	} else {
		f.AddValidationError("mcc__required", nil)
	}

	if body.Mnc != nil {
		if *body.Mnc >= 0 {
			f.Mnc = *body.Mnc
		} else {
			f.AddValidationError("mnc__invalid", nil)
		}
	} else {
		f.AddValidationError("mnc__required", nil)
	}

	if body.Price != nil {
		if *body.Price >= 0 {
			f.Price = model.Cost(*body.Price)
		} else {
			f.AddValidationError("price__invalid", nil)
		}
	} else {
		f.AddValidationError("price__required", nil)
	}

	if body.Currency != nil {
		currency, err := model.ParseCurrency(*body.Currency)
		if err != nil {
			f.AddValidationError("currency__invalid", err)
		}

		f.Currency = currency
	} else {
		f.AddValidationError("currency__required", nil)
	}

	return f
}

func (f priceForm) newPrice() model.Price {
	return model.NewPrice(
		f.AccountId,
		model.MccMnc{
			Mcc: f.Mcc,
			Mnc: f.Mnc,
		},
		f.Price,
		f.Currency)
}

type priceItemType string

var (
	priceItemGeneral = priceItemType("general")
	priceItemAccount = priceItemType("account")
)

type actionType string

var (
	actionTypeAdd    = actionType("add")
	actionTypeEdit   = actionType("edit")
	actionTypeRemove = actionType("delete")
)

type staffPriceResponse struct {
	*model.Price
	Enabled        bool          `json:"enabled"`
	PriceItemType  priceItemType `json:"type"`
	FormattedPrice string        `json:"price"`
	Actions        []actionType  `json:"actions"`
}

func newStaffPriceResponse(price model.Price, isGeneralPrice bool) staffPriceResponse {
	var piType priceItemType
	var actions []actionType
	if isGeneralPrice {
		piType = priceItemGeneral
		actions = append(actions, actionTypeEdit, actionTypeRemove)
	} else {
		if price.AccountId != nil {
			piType = priceItemAccount
			actions = append(actions, actionTypeEdit, actionTypeRemove)
		} else {
			piType = priceItemGeneral
			actions = append(actions, actionTypeAdd)
		}
	}

	return staffPriceResponse{
		Price:          &price,
		PriceItemType:  piType,
		Enabled:        price.Enabled,
		FormattedPrice: price.Price.Format(price.Currency),
		Actions:        actions,
	}
}

func newStaffPricesResponse(prices []model.Price, isGeneralPrice bool) []staffPriceResponse {
	pricesResp := make([]staffPriceResponse, len(prices))
	for idx, price := range prices {
		pricesResp[idx] = newStaffPriceResponse(price, isGeneralPrice)
	}

	return pricesResp
}

type updatePriceForm struct {
	*BaseForm
	Price   *model.Cost
	Enabled *bool `json:"enabled"`
}

func newUpdatePriceForm(r *http.Request) updatePriceForm {
	body := struct {
		Price   *int  `json:"price"`
		Enabled *bool `json:"enabled"`
	}{}
	f := updatePriceForm{
		BaseForm: &BaseForm{},
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		f.AddValidationError("body_json__invalid", err)
		return f
	}
	if body.Price != nil {
		if *body.Price >= 0 {
			newPrice := model.Cost(*body.Price)
			f.Price = &newPrice
		} else {
			f.AddValidationError("price__invalid", nil)
		}
	}
	if body.Enabled != nil {
		f.Enabled = body.Enabled
	}

	return f
}

func (f updatePriceForm) updatePriceItem(price *model.Price) {
	if f.Price != nil {
		price.Price = *f.Price
	}
	if f.Enabled != nil {
		price.Enabled = *f.Enabled
	}
}
