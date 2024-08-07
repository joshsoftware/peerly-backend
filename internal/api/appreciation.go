package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/peerly-backend/internal/app/appreciation"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/pkg/utils"
	logger "github.com/sirupsen/logrus"
)

func createAppreciationHandler(appreciationSvc appreciation.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var appreciation dto.Appreciation
		err := json.NewDecoder(req.Body).Decode(&appreciation)
		if err != nil {
			logger.Errorf("Error while decoding request data : %v", err)
			err = apperrors.JSONParsingErrorReq
			dto.ErrorRepsonse(rw, err)
			return
		}

		err = appreciation.ValidateCreateAppreciation()
		if err != nil {
			logger.Errorf("Error while validating request data : %v", err)
			dto.ErrorRepsonse(rw, err)
			return
		}

		resp, err := appreciationSvc.CreateAppreciation(req.Context(), appreciation)
		if err != nil {
			logger.Errorf("err : %v", err)
			dto.ErrorRepsonse(rw, err)
			return
		}
		dto.SuccessRepsonse(rw, http.StatusCreated, "Appreciation created successfully", resp)
	})
}

func getAppreciationByIDHandler(appreciationSvc appreciation.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		apprID, err := strconv.Atoi(vars["id"])
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}

		resp, err := appreciationSvc.GetAppreciationById(req.Context(), int32(apprID))
		if err != nil {
			logger.Errorf("err : %v", err)
			dto.ErrorRepsonse(rw, err)
			return
		}
		dto.SuccessRepsonse(rw, http.StatusOK, "Appreciation data got successfully", resp)
	})
}

// getAppreciationsHandler handles HTTP requests for appreciations
func listAppreciationsHandler(appreciationSvc appreciation.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		var filter dto.AppreciationFilter

		filter.Name = req.URL.Query().Get("name")
		filter.SortOrder = req.URL.Query().Get("sort_order")

		// Get pagination parameters
		page, limit := utils.GetPaginationParams(req)

		filter.Limit = limit
		filter.Page = page
		filter.Self = utils.GetSelfParam(req)

		appreciations, err := appreciationSvc.ListAppreciations(req.Context(), filter)
		if err != nil {
			logger.Errorf("err : %v", err)
			dto.ErrorRepsonse(rw, err)
			return
		}
		dto.SuccessRepsonse(rw, http.StatusOK, "Appreciations data got successfully ", appreciations)
	})
}

func deleteAppreciationHandler(appreciationSvc appreciation.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		apprId, err := strconv.Atoi(vars["id"])
		if err != nil {
			logger.Errorf("Error while decoding request param data : %v", err)
			dto.ErrorRepsonse(rw, apperrors.BadRequest)
			return
		}

		err = appreciationSvc.DeleteAppreciation(req.Context(), int32(apprId))
		if err != nil {
			logger.Errorf("err : %v", err)
			dto.ErrorRepsonse(rw, err)
			return
		}
		dto.SuccessRepsonse(rw, http.StatusOK, "Appreciation invalidate successfully", nil)
	})
}
