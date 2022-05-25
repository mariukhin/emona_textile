package api

import (
	"net/http"
)

type userCommon struct {
}

func (ac *userCommon) getPaginationOptions(w http.ResponseWriter, r *http.Request) {
	paginationOptions := []int{10, 25, 100}
	httpJsonResponse(w, paginationOptions)
}

type FilterOption struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

func NewFilterOption(name, title string) FilterOption {
	return FilterOption{
		Name:  name,
		Title: title,
	}
}

var (
	campaignFilterAll           = NewFilterOption("all", "All")
	campaignFilterScheduled     = NewFilterOption("scheduled", "Scheduled")
	campaignFilterInProgress    = NewFilterOption("in_progress", "In progress")
	campaignFilterFinished      = NewFilterOption("finished", "Finished")
	campaignFilterBlocked       = NewFilterOption("blocked", "Blocked")
	campaignFilterPaymentFailed = NewFilterOption("payment_failed", "Payment failed")
	campaignFilterCanceled      = NewFilterOption("canceled", "Canceled")
)

var campaignFilterOptions = []FilterOption{
	campaignFilterAll,
	campaignFilterScheduled,
	campaignFilterInProgress,
	campaignFilterFinished,
	campaignFilterBlocked,
	campaignFilterPaymentFailed,
	campaignFilterCanceled,
}

func (ac *userCommon) campaignFilters(w http.ResponseWriter, r *http.Request) {
	httpJsonResponse(w, campaignFilterOptions)
}
