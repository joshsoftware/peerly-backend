package repository

import (
	"context"
	"time"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

type OTPVerificationStorer interface {
	GetOTPVerificationStatus(ctx context.Context, otpReq dto.OTP) (otpInfo OTP, err error)
	GetCountOfOrgId(ctx context.Context, orgId int64) (count int, err error)
	ChangeIsVerifiedFlag(ctx context.Context, organizationID int64) error
	CreateOTPInfo(ctx context.Context, otpinfo OTP) error
	DeleteOTPData(ctx context.Context, orgId int64) error
}

type OTP struct {
	CreatedAt time.Time `db:"created_at"`
	OrgId     int64     `db:"org_id"`
	OTPCode   string    `db:"otp"`
}
