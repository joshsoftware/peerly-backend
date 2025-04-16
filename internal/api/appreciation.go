package api

import (
	"encoding/json"
	"fmt"
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
			log.Errorf(ctx, "Error while validating request data : %v", err)
			dto.ErrorRepsonse(rw, err)
			return
		}

		resp, err := appreciationSvc.CreateAppreciation(ctx, appreciation)
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
		resp, err := appreciationSvc.GetAppreciationById(ctx, int32(apprID))
		if err != nil {
			log.Errorf(ctx, "getAppreciationByIDHandler: err : %v", err)
			dto.ErrorRepsonse(rw, err)
			return
		}

		log.Debug(ctx, "getAppreciationByIDHandler: response: ", resp)
		log.Info(ctx, "Appreciation data fetched successfully")
		dto.SuccessRepsonse(rw, http.StatusOK, "Appreciation data fetched successfully", resp)
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
		appreciations, err := appreciationSvc.ListAppreciations(ctx, filter)
		if err != nil {
			log.Errorf(ctx, "listAppreciationsHandler: err : %v", err)
			dto.ErrorRepsonse(rw, err)
			return
		}
		log.Debug(ctx, "listAppreciationsHandler: response: ", appreciations)
		log.Info(ctx, "Appreciations data fetched successfully")
		dto.SuccessRepsonse(rw, http.StatusOK, "Appreciations data fetched successfully ", appreciations)
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
		err = appreciationSvc.DeleteAppreciation(ctx, int32(apprId))
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

func getAppreciationsHandler(apprSvc appreciation.Service, reportAppreciationSvc appreciation.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		qtrStr := req.URL.Query().Get("quarter")
		yearStr := req.URL.Query().Get("year")

		if qtrStr == "" || yearStr == "" {
			dto.ErrorRepsonse(rw, fmt.Errorf("quarter and year query parameters are required"))
			return
		}

		quarter, err := strconv.Atoi(qtrStr)
		if err != nil {
			dto.ErrorRepsonse(rw, fmt.Errorf("invalid quarter: %v", err))
			return
		}

		year, err := strconv.Atoi(yearStr)
		if err != nil {
			dto.ErrorRepsonse(rw, fmt.Errorf("invalid year: %v", err))
			return
		}
	
		if quarter < 1 || quarter > 4 {
			dto.ErrorRepsonse(rw, fmt.Errorf("quarter must be between 1 and 4"))
			return
		}


		tempFileName, err := apprSvc.GetAppreciations(ctx, quarter, year)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}

		rw.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", tempFileName))
		rw.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		http.ServeFile(rw, req, tempFileName)

	}
}