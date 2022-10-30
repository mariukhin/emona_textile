package api

import (
	"backend/app/logger"
	"backend/app/model"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

type carouselManagement struct {
	store     model.CarouselStore
	mediaHost string
}

func (cm *carouselManagement) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("w - %v", w)

	carouselItems, err := cm.store.FetchList()
	if err != nil {
		cm.getLog(ctx).Errorf("Fail to fetch carousel list - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	//if total == 0 {
	//	httpResponseError(w, http.StatusNotFound, "result__not_found")
	//	return
	//}

	resp := struct {
		// Total  int               `json:"total"`
		Result []*model.Carousel `json:"result"`
	}{
		// Total:  total,
		Result: carouselItems,
	}

	httpJsonResponse(w, resp)
}

func (cm *carouselManagement) getLog(ctx context.Context) *logrus.Entry {
	return logger.GetLogger(ctx).WithFields(logrus.Fields{
		"scope": "carousel_management",
	})
}
