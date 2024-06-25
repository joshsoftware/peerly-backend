package api

import (
	"encoding/json"
	"fmt"

	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/peerly-backend/internal/api/validations"
	"github.com/joshsoftware/peerly-backend/internal/app/organization"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"

	logger "github.com/sirupsen/logrus"
)

func listOrganizationHandler(orgSvc organization.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		organizations, err := orgSvc.ListOrganizations(req.Context())
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error listing organizations")
			apperrors.ErrorResp(rw, err)
			return
		}

		dto.Repsonse(rw, http.StatusOK, dto.SuccessResponse{Data: organizations})
	})
}

func getOrganizationHandler(orgSvc organization.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		fmt.Println("vars test: ",vars)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error id key is missing: request body conversion")
			apperrors.ErrorResp(rw, apperrors.InvalidId)
			return
		}

		organization, err := orgSvc.GetOrganization(req.Context(), id)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while fetching organization")
			apperrors.ErrorResp(rw, err)
			return
		}

		dto.Repsonse(rw, http.StatusOK, dto.SuccessResponse{Data: organization})
	})
}

func getOrganizationByDomainNameHandler(orgSvc organization.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		org, err := orgSvc.GetOrganizationByDomainName(req.Context(), vars["domainName"])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error retrieving organization by domain name: " + vars["domainName"])
			apperrors.ErrorResp(rw, err)
			return
		}
		dto.Repsonse(rw, http.StatusOK, dto.SuccessResponse{Data: org})
	})
}

func createOrganizationHandler(orgSvc organization.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		
		var organization dto.Organization
		err := json.NewDecoder(req.Body).Decode(&organization)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while decoding organization data")
			apperrors.ErrorResp(rw, apperrors.JSONParsingErrorReq)
			return
		}

		errorResponse, valid := validations.OrgValidate(organization)
		if !valid {
			respBytes, err := json.Marshal(errorResponse)
			if err != nil {
				logger.WithField("err", err.Error()).Error("Error marshaling organization data")
				// rw.WriteHeader(http.StatusInternalServerError)
				apperrors.ErrorResp(rw, apperrors.JSONParsingErrorReq)
				return
			}

			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(respBytes)
			return
		}
		organization.SubscriptionStatus = 1
		organization.CreatedBy = int64(req.Context().Value("roleid").(int))
		var createdOrganization dto.Organization
		createdOrganization, err = orgSvc.CreateOrganization(req.Context(), organization)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error create organization")
			apperrors.ErrorResp(rw, err)
			return
		}

		dto.Repsonse(rw, http.StatusCreated, dto.SuccessResponse{Data: createdOrganization})
	})
}

func updateOrganizationHandler(orgSvc organization.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		
		vars := mux.Vars(req)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error id key is missing")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		var organization dto.Organization
		err = json.NewDecoder(req.Body).Decode(&organization)
		if err != nil {
			apperrors.ErrorResp(rw, apperrors.JSONParsingErrorReq)
			return
		}
		organization.ID = int64(id)
		errorResponse, valid := validations.OrgUpdateValidate(organization)
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

		var updatedOrganization dto.Organization
		updatedOrganization, err = orgSvc.UpdateOrganization(req.Context(), organization)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while updating organization")
			apperrors.ErrorResp(rw, err)
			return
		}

		dto.Repsonse(rw, http.StatusOK, dto.SuccessResponse{Data: updatedOrganization})

	})
}

func deleteOrganizationHandler(orgSvc organization.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error id key is missing: request body conversion")
			apperrors.ErrorResp(rw, apperrors.InvalidId)
			return
		}

		userID, ok := req.Context().Value("userId").(int64)
		if !ok || userID == 0{
			http.Error(rw, "User ID not found in context", http.StatusInternalServerError)
			return
		}

		err = orgSvc.DeleteOrganization(req.Context(), id, userID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while deleting organization")
			apperrors.ErrorResp(rw, err)
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Header().Add("Content-Type", "application/json")
	})
}

func OTPVerificationHandler(orgSvc organization.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		
		var otpInfo dto.OTP
		err := json.NewDecoder(req.Body).Decode(&otpInfo)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while decoding otp data")
			apperrors.ErrorResp(rw, apperrors.JSONParsingErrorReq)
			return
		}

		errorResponse, valid := validations.OTPInfoValidate(otpInfo)
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
		fmt.Println("-------------------------------------------->")
		err = orgSvc.IsValidContactEmail(req.Context(), otpInfo)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while validating otp info")
			apperrors.ErrorResp(rw, err)
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Header().Add("Content-Type", "application/json")
	})
}

func ResendOTPhandler(orgSvc organization.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		orgId, err := strconv.Atoi(vars["id"])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error id key is missing: request body conversion")
			apperrors.ErrorResp(rw, apperrors.InvalidId)
			return
		}

		err = orgSvc.ResendOTPForContactEmail(req.Context(),int64(orgId))
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while resending otp ")
			apperrors.ErrorResp(rw, err)
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Header().Add("Content-Type", "application/json")


	})
}
