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
			dto.ErrorRepsonse(rw, err)
			return
		}

		dto.SuccessRepsonse(rw, http.StatusOK, "organization config fetched successfully",organization)
	})
}


func createOrganizationConfigHandler(orgSvc organizationConfig.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		
		var organization dto.OrganizationConfig
		err := json.NewDecoder(req.Body).Decode(&organization)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while decoding organization data")
			dto.ErrorRepsonse(rw, apperrors.JSONParsingErrorReq)
			return
		}

		err = validation.OrgValidate(organization)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}
		
		createdOrganizationConfig, err := orgSvc.CreateOrganizationConfig(req.Context(), organization)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error create organization")
			dto.ErrorRepsonse(rw, err)
			return
		}

		dto.SuccessRepsonse(rw, http.StatusCreated, "Organization Config Created Successfully" ,createdOrganizationConfig)
	})
}

func updateOrganizationConfigHandler(orgSvc organizationConfig.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		
		var organizationConfig dto.OrganizationConfig
		err := json.NewDecoder(req.Body).Decode(&organizationConfig)
		if err != nil {
			dto.ErrorRepsonse(rw, apperrors.JSONParsingErrorReq)
			return
		}
		organizationConfig.ID = 1
		err = validation.OrgUpdateValidate(organizationConfig)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}

		updatedOrganization, err := orgSvc.UpdateOrganizationConfig(req.Context(), organizationConfig)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while updating organization")
			dto.ErrorRepsonse(rw, err)
			return
		}

		dto.SuccessRepsonse(rw, http.StatusOK, "Organization Config Updated Successfully" ,updatedOrganization)

	})
}

