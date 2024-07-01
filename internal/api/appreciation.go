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
	logger "github.com/sirupsen/logrus"
)

func createAppreciationHandler(appreciationSvc appreciation.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var appreciation dto.Appreciation
		err := json.NewDecoder(req.Body).Decode(&appreciation)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while decoding request data")
			err = apperrors.JSONParsingErrorReq
			apperrors.ErrorResp(rw, err)
			return
		}

		errorResponse, ok := appreciation.CreateAppreciation()

		if !ok {
			respBytes, err := json.Marshal(errorResponse)
			if err != nil {
				logger.WithField("err", err.Error()).Error("Error marshaling organization data")
				apperrors.ErrorResp(rw, apperrors.JSONParsingErrorReq)
				return
			}
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(respBytes)
			return
		}
		resp, err := appreciationSvc.CreateAppreciation(req.Context(), appreciation)
		if err != nil {
			apperrors.ErrorResp(rw, err)
			return
		}
		dto.Repsonse(rw, http.StatusCreated, dto.SuccessResponse{Data: resp})
	})
}

func getAppreciationByIdHandler(appreciationSvc appreciation.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		apprId, err := strconv.Atoi(vars["id"])
		if err != nil {
			apperrors.ErrorResp(rw, apperrors.InvalidId)
			return
		}
		fmt.Println("appr: ",apprId)
		resp, err := appreciationSvc.GetAppreciationById(req.Context(), apprId)
		if err != nil {
			apperrors.ErrorResp(rw, err)
			return
		}
		dto.Repsonse(rw, http.StatusOK, dto.SuccessResponse{Data: resp})
	})
}

// getAppreciationsHandler handles HTTP requests for appreciations
func getAppreciationsHandler(appreciationSvc appreciation.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		var filter dto.AppreciationFilter

		// Extract query parameters or body fields
		filter.Name = req.URL.Query().Get("name")
		filter.SortOrder = req.URL.Query().Get("sort_order")

		// Call your appreciationService to fetch appreciations based on filter
		appreciations, err := appreciationSvc.GetAppreciation(req.Context(), filter)
		if err != nil {
			apperrors.ErrorResp(rw, err)
			return
		}
		dto.Repsonse(rw, http.StatusOK, dto.SuccessResponse{Data: appreciations})
	})
}

func validateAppreciationHandler(appreciationSvc appreciation.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		apprId, err := strconv.Atoi(vars["id"])
		if err != nil {
			apperrors.ErrorResp(rw, apperrors.BadRequest)
			return
		}
		
		res,err := appreciationSvc.ValidateAppreciation(req.Context(),false,apprId)
		if err != nil {
			apperrors.ErrorResp(rw, err)
			return 
		}
		if !res {
			apperrors.ErrorResp(rw, apperrors.InternalServer)
			return
		} 
		rw.WriteHeader(http.StatusOK)
	})
}