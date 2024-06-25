package organization

import (
	"context"
	"fmt"

	"time"

	email "github.com/joshsoftware/peerly-backend/internal/app/email"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/config"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/pkg/util"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	logger "github.com/sirupsen/logrus"
)

type service struct {
	OranizationRepo     repository.OrganizationStorer
	OTPVerificationRepo repository.OTPVerificationStorer
}

type Service interface {
	ListOrganizations(ctx context.Context) ([]dto.Organization, error)
	GetOrganization(ctx context.Context, id int) (dto.Organization, error)
	GetOrganizationByDomainName(ctx context.Context, domainName string) (dto.Organization, error)
	CreateOrganization(ctx context.Context, organization dto.Organization) (dto.Organization, error)
	UpdateOrganization(ctx context.Context, organization dto.Organization) (dto.Organization, error)
	DeleteOrganization(ctx context.Context, organizationID int, userId int64) (err error)
	IsValidContactEmail(ctx context.Context, otpInfo dto.OTP) (err error)
	ResendOTPForContactEmail(ctx context.Context, orgId int64) error
}

func NewService(oranizationRepo repository.OrganizationStorer, otpRepo repository.OTPVerificationStorer) Service {
	return &service{
		OranizationRepo:     oranizationRepo,
		OTPVerificationRepo: otpRepo,
	}
}

func (orgSvc *service) ListOrganizations(ctx context.Context) ([]dto.Organization, error) {

	organizations, err := orgSvc.OranizationRepo.ListOrganizations(ctx)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error listing organizations")
		return []dto.Organization{}, err
	}

	orgList := make([]dto.Organization, 0, len(organizations))
	for _, organization := range organizations {
		org := OrganizationDBToOrganization(organization)
		orgList = append(orgList, org)
	}

	return orgList, nil

}

func (orgSvc *service) GetOrganization(ctx context.Context, id int) (dto.Organization, error) {

	organization, err := orgSvc.OranizationRepo.GetOrganization(ctx, id)
	if err != nil {
		return dto.Organization{}, err
	}

	org := OrganizationDBToOrganization(organization)

	return org, nil

}

func (orgSvc *service) GetOrganizationByDomainName(ctx context.Context, domainName string) (dto.Organization, error) {

	organization, err := orgSvc.OranizationRepo.GetOrganizationByDomainName(ctx, domainName)
	if err != nil {
		return dto.Organization{}, err
	}
	org := OrganizationDBToOrganization(organization)
	return org, nil
}

func (orgSvc *service) CreateOrganization(ctx context.Context, organization dto.Organization) (dto.Organization, error) {

	isEmailPresent := orgSvc.OranizationRepo.IsEmailPresent(ctx, organization.ContactEmail)
	if isEmailPresent {
		return dto.Organization{}, apperrors.InvalidContactEmail
	}

	isDomainPresent := orgSvc.OranizationRepo.IsDomainPresent(ctx, organization.DomainName)
	if isDomainPresent {
		return dto.Organization{}, apperrors.InvalidDomainName
	}
	var createdOrganization repository.Organization
	createdOrganization, err := orgSvc.OranizationRepo.CreateOrganization(ctx, organization)
	if err != nil {
		return dto.Organization{}, err
	}
	org := OrganizationDBToOrganization(createdOrganization)
	OTPCode := util.GenerateRandomNumber(6)
	to := []string{org.ContactEmail}
	sub := "OTP Verification"
	
	mailData := &email.MailData{
		OTPCode: OTPCode,
	}
	mailReq := email.NewMail(config.ReadEnvString("SENDER_EMAIL"), to, sub, mailData)
	var otpinfo repository.OTP
	otpinfo.OTPCode = OTPCode
	otpinfo.OrgId = org.ID
	err = orgSvc.OTPVerificationRepo.CreateOTPInfo(ctx, otpinfo)
	if err != nil {
		return org, err
	}
	err = mailReq.SendMail()
	if err != nil {
		logger.Error("unable to send mail", "error", err)
		return org, nil
	}

	return org, nil
}

func (orgSvc *service) UpdateOrganization(ctx context.Context, organization dto.Organization) (dto.Organization, error) {

	if !orgSvc.OranizationRepo.IsOrganizationIdPresent(ctx, int64(organization.ID)) {
		return dto.Organization{}, apperrors.OrganizationNotFound
	}

	isEmailPresent := orgSvc.OranizationRepo.IsEmailPresent(ctx, organization.ContactEmail)
	if isEmailPresent {
		return dto.Organization{}, apperrors.InvalidContactEmail
	}

	isDomainPresent := orgSvc.OranizationRepo.IsDomainPresent(ctx, organization.DomainName)
	if isDomainPresent {
		return dto.Organization{}, apperrors.InvalidDomainName
	}

	updatedOrganization, err := orgSvc.OranizationRepo.UpdateOrganization(ctx, organization)
	if err != nil {
		return dto.Organization{}, err
	}

	org := OrganizationDBToOrganization(updatedOrganization)

	if organization.ContactEmail != "" {

		OTPCode := util.GenerateRandomNumber(6)
		to := []string{org.ContactEmail}
		sub := "OTP Verification"
		if err != nil {
			logger.WithField("err", err.Error()).Error("Database init failed")
			return org, nil
		}
		mailData := &email.MailData{
			OTPCode: OTPCode,
		}
		mailReq := email.NewMail(config.ReadEnvString("SENDER_EMAIL"), to, sub, mailData)
		var otpinfo repository.OTP
		otpinfo.OTPCode = OTPCode
		otpinfo.OrgId = org.ID
		err = orgSvc.OTPVerificationRepo.CreateOTPInfo(ctx, otpinfo)
		if err != nil {
			return org, err
		}
		err = mailReq.SendMail()
		if err != nil {
			logger.Error("unable to send mail", "error", err)
			return org, nil
		}

		return org, nil


		
	}
	return org, nil
}

func (orgSvc *service) DeleteOrganization(ctx context.Context, organizationID int, userId int64) (err error) {
	if !orgSvc.OranizationRepo.IsOrganizationIdPresent(ctx, int64(organizationID)) {
		return apperrors.OrganizationNotFound
	}

	err = orgSvc.OranizationRepo.DeleteOrganization(ctx, organizationID, userId)
	return err
}

func (orgSvc *service) IsValidContactEmail(ctx context.Context, otpInfo dto.OTP) error {
	fmt.Println("-------------------------------------------->serviceup")
	otp, err := orgSvc.OTPVerificationRepo.GetOTPVerificationStatus(ctx, otpInfo)
	if err != nil {
		if err == apperrors.InvalidReferenceId {
			logger.Error("otpInfo not found", "error", err)
			return err
		}
		logger.Error("unable to GetOTPVerificationStatus", "error", err)
		return err
	}

	fmt.Println("-------------------------------------------->service")
	fmt.Println("stored otp timestamp: ", otp.CreatedAt)
	fmt.Println("now: ", time.Now())
	expirationDuration := 2 * time.Minute
	expirationTime := otp.CreatedAt.Add(expirationDuration)
	if time.Now().After(expirationTime) {
		logger.Error("timelimit exceeded")
		return apperrors.TimeExceeded
	}

	if otp.OTPCode != otpInfo.OTPCode {
		logger.Error("invalid otp")
		return apperrors.InvalidOTP
	}

	err = orgSvc.OTPVerificationRepo.ChangeIsVerifiedFlag(ctx, otpInfo.OrgId)
	if err != nil {
		logger.Error("error from otp repo", "error", err)
		return err
	}

	err = orgSvc.OTPVerificationRepo.DeleteOTPData(ctx, otpInfo.OrgId)
	if err != nil {
		logger.Error("error in deleting otp data", "error", err)
		return err
	}

	return nil
}

func (orgSvc *service) ResendOTPForContactEmail(ctx context.Context, orgId int64) error {
	count, err := orgSvc.OTPVerificationRepo.GetCountOfOrgId(ctx, orgId)
	if err != nil {
		return err
	}
	if count == 0 {
		return apperrors.OrganizationNotFound
	}
	if count >= 3 {
		return apperrors.AttemptExceeded
	}

	org, err := orgSvc.OranizationRepo.GetOrganization(ctx, int(orgId))
	if err != nil {
		return err
	}
	OTPCode := util.GenerateRandomNumber(6)
	to := []string{org.ContactEmail}
	sub := "OTP Verification"
	mailData := &email.MailData{
		OTPCode: OTPCode,
	}
	mailReq := email.NewMail(config.ReadEnvString("SENDER_EMAIL"), to, sub, mailData)
	var otpinfo repository.OTP
	otpinfo.OTPCode = OTPCode
	otpinfo.OrgId = org.ID
	err = orgSvc.OTPVerificationRepo.CreateOTPInfo(ctx, otpinfo)
	if err != nil {
		return err
	}
	
	err = mailReq.SendMail()
	if err != nil {
		logger.Error("unable to send mail", "error", err)
		return apperrors.InernalServer
	}

	return nil
}
