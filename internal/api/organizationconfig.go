package api

import (
	"encoding/json"
	"net/http"

	"github.com/joshsoftware/peerly-backend/internal/app/organizationConfig"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	log "github.com/joshsoftware/peerly-backend/internal/pkg/logger"

	logger "github.com/sirupsen/logrus"
)



func getOrganizationConfigHandler(orgSvc organizationConfig.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		log.Debug(req.Context(),"getOrganizationConfigHandler: req: ",req)
		orgConfig, err := orgSvc.GetOrganizationConfig(req.Context())
		if err != nil {
			logger.Errorf("Error while fetching organization: %v",err)
			dto.ErrorRepsonse(rw, err)
			return
		}
		log.Debug(req.Context(),"getOrganizationConfigHandler: resp: ",orgConfig)
		dto.SuccessRepsonse(rw, http.StatusOK, "organization config fetched successfully",orgConfig)
	})
}


func createOrganizationConfigHandler(orgSvc organizationConfig.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		
		log.Debug(req.Context(),"createOrganizationConfigHandler: req: ",req)
		var orgConfig dto.OrganizationConfig
		err := json.NewDecoder(req.Body).Decode(&orgConfig)
		if err != nil {
			logger.Errorf("Error while decoding organization config data: %v",err)
			dto.ErrorRepsonse(rw, apperrors.JSONParsingErrorReq)
			return
		}

		err = orgConfig.OrgValidate()
		if err != nil {
			logger.Errorf("Error in validating request : %v",err)
			dto.ErrorRepsonse(rw, err)
			return
		}
		
		createdOrganizationConfig, err := orgSvc.CreateOrganizationConfig(req.Context(), orgConfig)
		if err != nil {
			logger.Errorf("Error in creating organization config: %v",err)
			dto.ErrorRepsonse(rw, err)
			return
		}
		log.Debug(req.Context(),"createOrganizationConfigHandler: resp: ",createdOrganizationConfig)
		dto.SuccessRepsonse(rw, http.StatusCreated, "Organization Config Created Successfully" ,createdOrganizationConfig)
	})
}

func updateOrganizationConfigHandler(orgSvc organizationConfig.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		
		var organizationConfig dto.OrganizationConfig
		err := json.NewDecoder(req.Body).Decode(&organizationConfig)
		if err != nil {
			logger.Errorf("Error while decoding organization data: %v",err)
			dto.ErrorRepsonse(rw, apperrors.JSONParsingErrorReq)
			return
		}
		
		log.Debug(req.Context(),"updateOrganizationConfigHandler: request: ",req)
		organizationConfig.ID = 1
		err = organizationConfig.OrgUpdateValidate()
		if err != nil {
			logger.Errorf("Error in validating request : %v",err)
			dto.ErrorRepsonse(rw, err)
			return
		}

		updatedOrganization, err := orgSvc.UpdateOrganizationConfig(req.Context(), organizationConfig)
		if err != nil {
			logger.Errorf("Error while updating organization: %v",err)
			dto.ErrorRepsonse(rw, err)
			return
		}

		log.Debug(req.Context(),"updateOrganizationConfigHandler: resp: ",updatedOrganization)
		dto.SuccessRepsonse(rw, http.StatusOK, "Organization Config Updated Successfully" ,updatedOrganization)

	})
}

