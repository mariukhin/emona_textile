package api

import (
	"amifactory.team/sequel/coton-app-backend/app/logger"
	"amifactory.team/sequel/coton-app-backend/app/model"
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
)

type staffRole struct {
	store model.StaffStore
}

func (s *staffRole) getStaffRoles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	staff, err := model.StaffFromContext(r.Context(), model.KCtxKeyStaffMe)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "request__unauthorized")
		return
	}

	// TODO add 'Admin' role check
	if staff == nil {
		// TODO
	}

	roles, err := s.store.FetchStaffRoles()
	if err != nil {
		s.getLog(ctx).Errorf("Fail to fetch staff roles - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	httpJsonResponse(w, roles)
}

func (s *staffRole) getLog(ctx context.Context) *logrus.Entry {
	return logger.GetLogger(ctx).WithFields(logrus.Fields{
		"scope": "staff_role",
	})
}
