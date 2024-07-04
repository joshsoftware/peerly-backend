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
			dto.ErrorRepsonse(rw, err,nil)
			return
		}

		errorResponse, ok := appreciation.CreateAppreciation()

		if !ok {
			dto.ErrorRepsonse(rw, apperrors.BadRequest,errorResponse)
			return
		}
		resp, err := appreciationSvc.CreateAppreciation(req.Context(), appreciation)
		if err != nil {
			dto.ErrorRepsonse(rw, err,nil)
			return
		}
		dto.SuccessRepsonse(rw, http.StatusCreated,"Appreciation created successfully" ,resp)
	})
}

func getAppreciationByIdHandler(appreciationSvc appreciation.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		apprId, err := strconv.Atoi(vars["id"])
		if err != nil {
			dto.ErrorRepsonse(rw, err,nil)
			return
		}
		fmt.Println("appr: ",apprId)
		resp, err := appreciationSvc.GetAppreciationById(req.Context(), apprId)
		if err != nil {
			dto.ErrorRepsonse(rw, err,nil)
			return
		}
		dto.SuccessRepsonse(rw, http.StatusOK,"Appreciation data got successfully" , resp)
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
			dto.ErrorRepsonse(rw, err,nil)
			return
		}
		dto.SuccessRepsonse(rw, http.StatusOK,"Appreciations data got successfully " ,appreciations)
	})
}

func validateAppreciationHandler(appreciationSvc appreciation.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		apprId, err := strconv.Atoi(vars["id"])
		if err != nil {
			dto.ErrorRepsonse(rw, apperrors.BadRequest,nil)
			return
		}
		
		res,err := appreciationSvc.ValidateAppreciation(req.Context(),false,apprId)
		if err != nil {
			dto.ErrorRepsonse(rw, err,nil)
			return 
		}
		if !res {
			dto.ErrorRepsonse(rw, apperrors.InternalServer,nil)
			return
		} 
		dto.SuccessRepsonse(rw,http.StatusOK,"Appreciation invalidate successfully",nil)
	})
}