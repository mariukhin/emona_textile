package model

import (
	"amifactory.team/sequel/coton-app-backend/app/logger"
	"context"
	"fmt"
	"time"
)

type CampaignTarget string

var (
	Manual = CampaignTarget("manual")
	Group  = CampaignTarget("group")
)

var campaignTargets = []CampaignTarget{
	Manual, Group,
}

func ParseCampaignTarget(value string) (CampaignTarget, error) {
	for _, target := range campaignTargets {
		if target == CampaignTarget(value) {
			return target, nil
		}
	}

	return "unknown", fmt.Errorf("fail to parse CampaignTarget - '%s' unknown value", value)
}

type SendOption string

var (
	Immediately = SendOption("immediately")
	Scheduled   = SendOption("scheduled")
)

var sendOptions = []SendOption{
	Immediately, Scheduled,
}

func ParseSendOption(value string) (SendOption, error) {
	for _, option := range sendOptions {
		if option == SendOption(value) {
			return option, nil
		}
	}

	return "unknown", fmt.Errorf("fail to parse SendOption - '%s' unknown value", value)
}

type CampaignRoute struct {
	ID            string `bson:"_id" json:"-"`
	CampaignID    string `bson:"campaign" json:"-"`
	Country       string `bson:"country" json:"country"`
	MccMnc        `bson:",inline"`
	MessagesCount int      `bson:"messages_count" json:"-"`
	SmsCount      int      `bson:"sms_count" json:"sms_count"`
	Currency      Currency `bson:"currency,omitempty" json:"currency,omitempty"`
	Price         Cost     `bson:"price" json:"-"`
	Total         Cost     `bson:"total" json:"-"`
	Available     bool     `bson:"available" json:"available"`
}

type Campaign struct {
	ID string `bson:"_id" json:"id"`

	CreatedAt time.Time `bson:"cat" json:"-"`
	UpdatedAt time.Time `bson:"uat" json:"-"`

	AccountID       string `bson:"account" json:"-"`
	AccountMemberID string `bson:"member" json:"-"`

	Name        string `bson:"name" json:"name"`
	SenderID    string `bson:"sender_id" json:"sender_id"`
	MessageText string `bson:"message_text" json:"message_text"`

	Target  CampaignTarget `bson:"targets_type" json:"targets_type"`
	GroupID string         `bson:"group,omitempty" json:"group,omitempty"`
	Phones  []*Phone       `bson:"phones,omitempty" json:"phones,omitempty"`

	SendOption       SendOption `bson:"send_option" json:"send_option"`
	SendAfter        *time.Time `bson:"send_after,omitempty" json:"send_after,omitempty"`
	TrackLinkOpening bool       `bson:"track_link_opening" json:"track_link_opening"`

	PriceCurrency Currency `bson:"currency" json:"-"`

	MessageParts    int `bson:"message_parts" json:"message_parts"`
	MessagePartsMax int `bson:"message_parts_max,omitempty" json:"message_parts_max,omitempty"`
	MessagesCount   int `bson:"messages_count" json:"messages_count"`
	SmsCount        int `bson:"sms_count" json:"sms_count"`

	Currency Currency `bson:"currency" json:"-"`
	Total    Cost     `bson:"sms_count" json:"-"`

	Routes []CampaignRoute `bson:"routes" json:"-"`
	// TODO other counters
}

func (c Campaign) GroupPhones() (map[MccMnc]int, error) {
	if c.Target != Manual {
		return nil, fmt.Errorf("fail to group phones - wrong campaign target type %v", c.Target)
	}

	phoneGroups := make(map[MccMnc]int)
	for _, phone := range c.Phones {
		count := phoneGroups[phone.MccMnc]
		phoneGroups[phone.MccMnc] = count + 1
	}

	return phoneGroups, nil
}

func (c *Campaign) SetRoutes(routes []CampaignRoute) {
	smsTotal := 0
	messagesTotal := 0
	var total Cost

	for _, route := range routes {
		if route.Available {
			smsTotal += route.SmsCount
			messagesTotal += route.MessagesCount
			total = total.Add(route.Total)
		}
	}

	c.Routes = routes
	c.SmsCount = smsTotal
	c.MessagesCount = messagesTotal
	c.Total = total
}

type CampaignStore interface {
	EnsureIndexes(ctx context.Context) error

	Add(ctx context.Context, campaign *Campaign) error
}

func NewCampaignStore(store Storage) (CampaignStore, error) {
	return &campaignStore{
		store: store,
	}, nil
}

type campaignStore struct {
	store Storage
}

func (cs *campaignStore) EnsureIndexes(ctx context.Context) error {
	log := logger.GetLogger(ctx).WithField("store", "campaign")

	log.Info("Fetching storage indexes")

	log.Info("Storage indexes are up-to-date")

	return nil
}

func (cs *campaignStore) Add(ctx context.Context, campaign *Campaign) error {
	_, err := cs.store.Collection("campaigns").InsertOne(ctx, campaign)
	if err != nil {
		return fmt.Errorf("fail to add campaign - %v", err)
	}

	return nil
}
