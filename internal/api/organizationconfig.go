package api

import (
	"encoding/json"
	"net/http"

	"github.com/joshsoftware/peerly-backend/internal/app/organizationConfig"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
)



func getOrganizationConfigHandler(orgSvc organizationConfig.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		logger.Debug(req.Context(),"getOrganizationConfigHandler: req: ",req)
		orgConfig, err := orgSvc.GetOrganizationConfig(req.Context())
		if err != nil {
			logger.Errorf(req.Context(),"Error while fetching organization: %v",err)
			dto.ErrorRepsonse(rw, err)
			return
		}
		logger.Debug(req.Context(),"getOrganizationConfigHandler: resp: ",orgConfig)
		logger.Info(req.Context(),"organization config fetched successfully")
		dto.SuccessRepsonse(rw, http.StatusOK, "organization config fetched successfully",orgConfig)
	})
}


func createOrganizationConfigHandler(orgSvc organizationConfig.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		
		logger.Debug(req.Context(),"createOrganizationConfigHandler: request: ",req)
		var orgConfig dto.OrganizationConfig
		err := json.NewDecoder(req.Body).Decode(&orgConfig)
		if err != nil {
			logger.Errorf(req.Context(),"Error while decoding organization config data: %v",err)
			dto.ErrorRepsonse(rw, apperrors.JSONParsingErrorReq)
			return
		}

		err = orgConfig.OrgValidate()
		if err != nil {
			logger.Errorf(req.Context(),"Error in validating request : %v",err)
			dto.ErrorRepsonse(rw, err)
			return
		}
		
		createdOrganizationConfig, err := orgSvc.CreateOrganizationConfig(req.Context(), orgConfig)
		if err != nil {
			logger.Errorf(req.Context(),"Error in creating organization config: %v",err)
			dto.ErrorRepsonse(rw, err)
			return
		}
		logger.Debug(req.Context(),"createOrganizationConfigHandler: resp: ",createdOrganizationConfig)
		logger.Info(req.Context(),"Organization Config Created Successfully")
		dto.SuccessRepsonse(rw, http.StatusCreated, "Organization Config Created Successfully" ,createdOrganizationConfig)
	})
}

func updateOrganizationConfigHandler(orgSvc organizationConfig.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		
		var organizationConfig dto.OrganizationConfig
		err := json.NewDecoder(req.Body).Decode(&organizationConfig)
		if err != nil {
			logger.Errorf(req.Context(),"Error while decoding organization data: %v",err)
			dto.ErrorRepsonse(rw, apperrors.JSONParsingErrorReq)
			return
		}
		
		logger.Info(req.Context(),"updateOrganizationConfigHandler: request: ",req)
		organizationConfig.ID = 1
		err = organizationConfig.OrgUpdateValidate()
		if err != nil {
			logger.Errorf(req.Context(),"Error in validating request : %v",err)
			dto.ErrorRepsonse(rw, err)
			return
		}

		updatedOrganization, err := orgSvc.UpdateOrganizationConfig(req.Context(), organizationConfig)
		if err != nil {
			logger.Errorf(req.Context(),"Error while updating organization: %v",err)
			dto.ErrorRepsonse(rw, err)
			return
		}

		logger.Debug(req.Context(),"updateOrganizationConfigHandler: resp: ",updatedOrganization)
		logger.Info(req.Context(),"Organization Config Updated Successfully" )
		dto.SuccessRepsonse(rw, http.StatusOK, "Organization Config Updated Successfully" ,updatedOrganization)

	})
}

