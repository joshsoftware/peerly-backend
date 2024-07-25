package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/peerly-backend/internal/app/badges"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/sirupsen/logrus"
)

func createBadgeHandler(badgeSvc badges.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var badge dto.Badge
		err := json.NewDecoder(req.Body).Decode(&badge)
		if err != nil {
			logger.Errorf("error while decoding request data, err: %s", err.Error())
			err = apperrors.JSONParsingErrorReq
			dto.ErrorRepsonse(rw, err)
			return
		}

		err = badge.ValidateCreateBadge()
		if err != nil {
			logger.Error(fmt.Sprintf("invalid badge request: %v", err))
			dto.ErrorRepsonse(rw, err)
			return
		}

		resp, err := badgeSvc.CreateBadge(req.Context(), badge)
		if err != nil {
			logger.Error(fmt.Sprintf("err: %v", err))
			dto.ErrorRepsonse(rw, err)
			return
		}

		dto.SuccessRepsonse(rw, http.StatusCreated, "Badge created successfully", resp)
	})
}

func listBadgesHandler(badgeSvc badges.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		resp, err := badgeSvc.ListBadges(req.Context())
		if err != nil {
			logger.Errorf("err: %v", err)
			dto.ErrorRepsonse(rw, err)
			return
		}
		dto.SuccessRepsonse(rw, http.StatusOK, "List fetched successfully", resp)
	})
}

func getBadgeHandler(badgeSvc badges.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		badgeID, err := strconv.Atoi(vars["id"])
		if err != nil {
			logger.Error(fmt.Sprintf("Error while decoding request param data : %v", err))
			dto.ErrorRepsonse(rw, apperrors.BadRequest)
			return
		}

		resp, err := badgeSvc.GetBadge(req.Context(),int8(badgeID))
		if err != nil {
			logger.Errorf("err: %v", err)
			dto.ErrorRepsonse(rw, err)
			return
		}
		dto.SuccessRepsonse(rw, http.StatusOK, "badge fetched successfully", resp)

	})
}

func deleteBadgeHandler(badgeSvc badges.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		badgeID, err := strconv.Atoi(vars["id"])
		if err != nil {
			logger.Error(fmt.Sprintf("Error while decoding request badge id : %v", err))
			dto.ErrorRepsonse(rw, apperrors.BadRequest)
			return
		}

		err = badgeSvc.DeleteBadge(req.Context(),int8(badgeID))
		if err != nil {
			logger.Errorf("err: %v", err)
			dto.ErrorRepsonse(rw, err)
			return
		}
		dto.SuccessRepsonse(rw, http.StatusOK, "badge deleted successfully", nil)

	})
}

func updateBadgeHandler(badgeSvc badges.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		badgeID, err := strconv.Atoi(vars["id"])
		if err != nil {
			logger.Error(fmt.Sprintf("Error while decoding request badge id : %v", err))
			dto.ErrorRepsonse(rw, apperrors.BadRequest)
			return
		}

		var badge dto.Badge
		err = json.NewDecoder(req.Body).Decode(&badge)
		if err != nil {
			logger.Errorf("error while decoding request data, err: %s", err.Error())
			err = apperrors.JSONParsingErrorReq
			dto.ErrorRepsonse(rw, err)
			return
		}

		badge.ID = int8(badgeID)
		err = badge.ValidateUpdateBadge()
		if err != nil {
			logger.Error(fmt.Sprintf("invalid badge request: %v", err))
			dto.ErrorRepsonse(rw, err)
			return
		}

		resp, err := badgeSvc.UpdateBadge(req.Context(), badge)
		if err != nil {
			logger.Errorf("err: %v", err)
			dto.ErrorRepsonse(rw, err)
			return
		}

		dto.SuccessRepsonse(rw, http.StatusOK, "Badge updated successfully", resp)
	})
}
