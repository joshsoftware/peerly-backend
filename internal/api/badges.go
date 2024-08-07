package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/peerly-backend/internal/app/badges"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/sirupsen/logrus"
)

func listBadgesHandler(badgeSvc badges.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		resp, err := badgeSvc.ListBadges(ctx)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}
		dto.SuccessRepsonse(rw, 200, "badges list fetched successfully", resp)
	})
}

func editBadgesHandler(badgeSvc badges.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		vars := mux.Vars(req)
		var reqData dto.UpdateBadgeReq
		err := json.NewDecoder(req.Body).Decode(&reqData)
		if err != nil {
			logger.Errorf("error while decoding request data, err: %s", err.Error())
			err = apperrors.JSONParsingErrorReq
			dto.ErrorRepsonse(rw, err)
			return
		}
		err = badgeSvc.EditBadge(ctx, vars["id"], reqData.RewardPoints)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}
		dto.SuccessRepsonse(rw, 200, "badge reward points updated successfully", nil)
	})
}
