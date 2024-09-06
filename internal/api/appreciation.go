package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/peerly-backend/internal/app/appreciation"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	log "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
	"github.com/joshsoftware/peerly-backend/internal/pkg/utils"
)

func createAppreciationHandler(appreciationSvc appreciation.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var appreciation dto.Appreciation
		ctx := req.Context()
		err := json.NewDecoder(req.Body).Decode(&appreciation)
		if err != nil {
			log.Errorf(ctx, "Error while decoding request data : %v", err)
			err = apperrors.JSONParsingErrorReq
			dto.ErrorRepsonse(rw, err)
			return
		}

		log.Debug(ctx, "createAppreciationHandler: request: ", req)
		err = appreciation.ValidateCreateAppreciation()
		if err != nil {
			log.Errorf(req.Context(), "Error while validating request data : %v", err)
			dto.ErrorRepsonse(rw, err)
			return
		}

		resp, err := appreciationSvc.CreateAppreciation(req.Context(), appreciation)
		if err != nil {
			log.Errorf(ctx, "createAppreciationHandler: err : %v", err)
			dto.ErrorRepsonse(rw, err)
			return
		}
		log.Debug(ctx, "createAppreciationHandler: response: ", resp)
		log.Info(ctx, "Appreciation created successfully")
		dto.SuccessRepsonse(rw, http.StatusCreated, "Appreciation created successfully", resp)
	})
}

func getAppreciationByIDHandler(appreciationSvc appreciation.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		vars := mux.Vars(req)
		apprID, err := strconv.Atoi(vars["id"])
		if err != nil {
			log.Errorf(ctx, "Error while decoding appreciation id : %v", err)
			dto.ErrorRepsonse(rw, err)
			return
		}

		log.Debug(ctx, "getAppreciationByIDHandler: request: ", req)
		resp, err := appreciationSvc.GetAppreciationById(req.Context(), int32(apprID))
		if err != nil {
			log.Errorf(ctx, "getAppreciationByIDHandler: err : %v", err)
			dto.ErrorRepsonse(rw, err)
			return
		}

		log.Debug(ctx, "getAppreciationByIDHandler: response: ", resp)
		log.Info(ctx, "Appreciation data fetch successfully")
		dto.SuccessRepsonse(rw, http.StatusOK, "Appreciation data fetch successfully", resp)
	})
}

// getAppreciationsHandler handles HTTP requests for appreciations
func listAppreciationsHandler(appreciationSvc appreciation.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var filter dto.AppreciationFilter
		ctx := req.Context()
		filter.Name = req.URL.Query().Get("name")
		filter.SortOrder = req.URL.Query().Get("sort_order")

		// Get pagination parameters
		page, limit := utils.GetPaginationParams(req)

		filter.Limit = limit
		filter.Page = page
		filter.Self = utils.GetSelfParam(req)
		log.Debug(ctx, "listAppreciationsHandler: request: ", req)
		appreciations, err := appreciationSvc.ListAppreciations(req.Context(), filter)
		if err != nil {
			log.Errorf(req.Context(), "listAppreciationsHandler: err : %v", err)
			dto.ErrorRepsonse(rw, err)
			return
		}
		log.Debug(ctx, "listAppreciationsHandler: response: ", appreciations)
		log.Info(ctx, "Appreciations data fetch successfully")
		dto.SuccessRepsonse(rw, http.StatusOK, "Appreciations data fetch successfully ", appreciations)
	})
}

func deleteAppreciationHandler(appreciationSvc appreciation.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		ctx := req.Context()
		apprId, err := strconv.Atoi(vars["id"])
		if err != nil {
			log.Errorf(ctx, "Error while decoding request param data : %v", err)
			dto.ErrorRepsonse(rw, apperrors.BadRequest)
			return
		}

		log.Info(ctx, "deleteAppreciationHandler: request: ", req)
		err = appreciationSvc.DeleteAppreciation(req.Context(), int32(apprId))
		if err != nil {
			log.Errorf(ctx, "deleteAppreciationHandler: err : %v", err)
			dto.ErrorRepsonse(rw, err)
			return
		}
		log.Debug(ctx, "deleteAppreciationHandler: resp: ", err)
		log.Info(ctx, "Appreciation invalidate successfully")
		dto.SuccessRepsonse(rw, http.StatusOK, "Appreciation invalidate successfully", nil)
	})
}
