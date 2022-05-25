package api

import (
	"amifactory.team/sequel/coton-app-backend/app/logger"
	"amifactory.team/sequel/coton-app-backend/app/model"
	"context"
	"errors"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"math"
	"net/http"
	"regexp"
	"strconv"
)

type usersManagement struct {
	store     model.UserStore
	mediaHost string
}

func (u *usersManagement) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	srcParam := r.URL.Query().Get("src")
	pageParam := r.URL.Query().Get("page")
	pageSizeParam := r.URL.Query().Get("page_size")

	page := 0
	if len(pageParam) > 0 {
		page, err := strconv.Atoi(pageParam)
		if err != nil {
			u.getLog(ctx).Errorf("Fail to get parse 'page' param - %v", err)
			httpResponseError(w, http.StatusBadRequest, "page__invalid")
			return
		}

		if page < 0 {
			httpResponseError(w, http.StatusBadRequest, "page__out_of_bounds")
			return
		}
	}

	pageSize := 10
	if len(pageSizeParam) > 0 {
		pageSize, err := strconv.Atoi(pageSizeParam)
		if err != nil {
			u.getLog(ctx).Errorf("Fail to get parse 'page_size' param - %v", err)
			httpResponseError(w, http.StatusBadRequest, "page_size__invalid")
			return
		}

		if pageSize <= 0 || 100 < pageSize {
			httpResponseError(w, http.StatusBadRequest, "page_size__out_of_bounds")
			return
		}
	}

	var searchPhase *string
	if len(srcParam) > 0 {
		if len(srcParam) >= 3 {
			srcEscaped := regexp.QuoteMeta(srcParam)
			searchPhase = &srcEscaped
		} else {
			httpResponseError(w, http.StatusBadRequest, "src__invalid")
			return
		}
	} else {
		searchPhase = nil
	}

	users, total, err := u.store.FetchList(page*pageSize, pageSize, searchPhase)
	if err != nil {
		if errors.Is(err, model.ErrPageOutOfBounds) {
			httpResponseError(w, http.StatusBadRequest, "page__out_of_bounds")
			return
		}

		u.getLog(ctx).Errorf("Fail to fetch users list - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	if total == 0 && (searchPhase != nil) {
		httpResponseError(w, http.StatusNotFound, "result__not_found")
		return
	}

	resp := struct {
		Pages  int                `json:"pages"`
		Total  int                `json:"total"`
		Result []*model.UserShort `json:"result"`
	}{
		Pages:  int(math.Ceil(float64(total) / float64(pageSize))),
		Total:  total,
		Result: users,
	}

	httpJsonResponse(w, resp)
}

func (u *usersManagement) fetch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId := chi.URLParam(r, "userID")
	if len(userId) == 0 {
		httpResponseError(w, http.StatusBadRequest, "user__invalid")
		return
	}

	user, err := u.store.FindByID(r.Context(), userId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			httpResponseError(w, http.StatusNotFound, "user__not_found")
		} else {
			u.getLog(ctx).Errorf("Fail to fetch user by ID - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
		}

		return
	}

	user.UpdatePhotoURL(u.mediaHost)
	httpJsonResponse(w, user)
}

func (u *usersManagement) getLog(ctx context.Context) *logrus.Entry {
	return logger.GetLogger(ctx).WithFields(logrus.Fields{
		"scope": "user_management",
	})
}
