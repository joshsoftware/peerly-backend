package api

import (
	"encoding/json"
	"net/http"

	"github.com/joshsoftware/peerly-backend/internal/api/validation"
	"github.com/joshsoftware/peerly-backend/internal/app/organizationConfig"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"

	logger "github.com/sirupsen/logrus"
)



func getOrganizationConfigHandler(orgSvc organizationConfig.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		organization, err := orgSvc.GetOrganizationConfig(req.Context())
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while fetching organization")
			apperrors.ErrorResp(rw, err)
			return
		}

		dto.Repsonse(rw, http.StatusOK, dto.SuccessResponse{Data: organization})
	})
}


func createOrganizationConfigHandler(orgSvc organizationConfig.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		
		var organization dto.OrganizationConfig
		err := json.NewDecoder(req.Body).Decode(&organization)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while decoding organization data")
			apperrors.ErrorResp(rw, apperrors.JSONParsingErrorReq)
			return
		}

		errorResponse, valid := validation.OrgValidate(organization)
		if !valid {
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
		
		createdOrganizationConfig, err := orgSvc.CreateOrganizationConfig(req.Context(), organization)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error create organization")
			apperrors.ErrorResp(rw, err)
			return
		}

		dto.Repsonse(rw, http.StatusCreated, dto.SuccessResponse{Data: createdOrganizationConfig})
	})
}

func updateOrganizationConfigHandler(orgSvc organizationConfig.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		
		var organizationConfig dto.OrganizationConfig
		err := json.NewDecoder(req.Body).Decode(&organizationConfig)
		if err != nil {
			apperrors.ErrorResp(rw, apperrors.JSONParsingErrorReq)
			return
		}
		organizationConfig.ID = 1
		errorResponse, valid := validation.OrgUpdateValidate(organizationConfig)
		if !valid {
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
		updatedOrganization, err := orgSvc.UpdateOrganizationConfig(req.Context(), organizationConfig)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while updating organization")
			apperrors.ErrorResp(rw, err)
			return
		}

		dto.Repsonse(rw, http.StatusOK, dto.SuccessResponse{Data: updatedOrganization})

	})
}

