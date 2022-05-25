package api

import (
	"context"
	"errors"
	"net/http"

	"amifactory.team/sequel/coton-app-backend/app/logger"
	"amifactory.team/sequel/coton-app-backend/app/model"
	"github.com/sirupsen/logrus"
)

type userPrices struct {
	pricesStore  model.PricesStore
	accountStore model.AccountStore
	phoneService model.PhoneService
}

func (p *userPrices) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	account, _ := model.AccountFromContext(ctx)

	form := NewPageForm(r, 1)
	if form.HasValidationErrors() {
		httpResponseErrors(w, http.StatusBadRequest, form.ValidationErrors)
		return
	}

	filter := model.PricesFilter{
		SearchPhase: form.SearchPhrase,
		AccountID:   &account.ID,
		Currency:    &account.Billing.Currency,
	}.SetEnabled(true)

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

	resp := newUserPricesResponse(prices)
	httpJsonResponse(w, NewPageResp(total, form.PageSize, resp))
}

func (p *userPrices) getLog(ctx context.Context) *logrus.Entry {
	return logger.GetLogger(ctx).WithFields(logrus.Fields{
		"scope": "user-prices",
	})
}

type userPriceResponse struct {
	*model.Price
	FormattedPrice string `json:"price"`
	LastUpdate     string `json:"last_update"`
}

func newUserPriceResponse(price model.Price) userPriceResponse {
	return userPriceResponse{
		Price:          &price,
		FormattedPrice: price.Price.Format(price.Currency),
		LastUpdate:     price.UpdatedAt.Format("01/02/2006 15:04"),
	}
}

func newUserPricesResponse(prices []model.Price) []userPriceResponse {
	pricesResp := make([]userPriceResponse, len(prices))
	for idx, price := range prices {
		pricesResp[idx] = newUserPriceResponse(price)
	}

	return pricesResp
}
