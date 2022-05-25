package api

import (
	"amifactory.team/sequel/coton-app-backend/app/model"
	"net/http"
)

type adminCommon struct {
}

type AccountStatusFilter struct {
	value         string
	accountStatus *model.AccountStatus
}

func (f AccountStatusFilter) MarshalText() ([]byte, error) {
	return []byte(f.value), nil
}

func NewAccountStatusFilter(status model.AccountStatus) AccountStatusFilter {
	if status.IsValid() {
		return AccountStatusFilter{
			value:         string(status),
			accountStatus: &status,
		}
	} else {
		return AccountStatusFilter{
			value:         string(status),
			accountStatus: nil,
		}
	}
}

var (
	AccountStatusFilterAll     = NewAccountStatusFilter("all")
	AccountStatusFilterCreated = NewAccountStatusFilter(model.AccountStatusCreated)
	AccountStatusFilterActive  = NewAccountStatusFilter(model.AccountStatusActive)
	AccountStatusFilterBlocked = NewAccountStatusFilter(model.AccountStatusBlocked)
)

var accountStatusFilters = []AccountStatusFilter{
	AccountStatusFilterAll,
	AccountStatusFilterCreated,
	AccountStatusFilterActive,
	AccountStatusFilterBlocked,
}

func ParseAccountStatusFilter(val string) *AccountStatusFilter {
	for _, filter := range accountStatusFilters {
		if val == filter.value {
			return &filter
		}
	}

	return nil
}

func (ac *adminCommon) getPaginationOptions(w http.ResponseWriter, r *http.Request) {
	paginationOptions := []int{10, 25, 100}
	httpJsonResponse(w, paginationOptions)
}

func (ac *adminCommon) getAccountStatus(w http.ResponseWriter, r *http.Request) {
	httpJsonResponse(w, accountStatusFilters)
}

func (ac *adminCommon) getModerationOptions(w http.ResponseWriter, r *http.Request) {
	options := []model.AccountModeration{
		model.AccountModerationOn,
		model.AccountModerationOff,
	}
	httpJsonResponse(w, options)
}

func (ac *adminCommon) getAccountMemberRoles(w http.ResponseWriter, r *http.Request) {
	type role struct {
		Name  string `json:"name"`
		Title string `json:"title"`
	}

	roles := []role{
		role{
			Name:  string(model.AccountMemberRoleOwner),
			Title: "Owner",
		},
		role{
			Name:  string(model.AccountMemberRoleAdmin),
			Title: "Admin",
		},
		role{
			Name:  string(model.AccountMemberRoleUser),
			Title: "User",
		},
	}

	httpJsonResponse(w, roles)
}

func (ac *adminCommon) getCurrencies(w http.ResponseWriter, r *http.Request) {
	type currencyResp struct {
		Name  string `json:"name"`
		Title string `json:"title"`
	}
	currencies := []currencyResp{
		currencyResp{
			Name:  string(model.Euro),
			Title: "€",
		},
	}
	httpJsonResponse(w, currencies)
}
