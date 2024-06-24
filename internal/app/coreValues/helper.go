package corevalues

import (
	"context"
	"fmt"
	"strconv"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	logger "github.com/sirupsen/logrus"
)

func validateParentCoreValue(ctx context.Context, storer repository.CoreValueStorer, organisationID, coreValueID int64) (ok bool) {
	coreValue, err := storer.GetCoreValue(ctx, organisationID, coreValueID)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Parent core value id not present")
		return
	}

	if coreValue.ParentID != nil || coreValue.SoftDelete {
		logger.Error("Invalid parent core value id")
		return
	}

	return true
}

func Validate(ctx context.Context, coreValue dto.CreateCoreValueReq, storer repository.CoreValueStorer, organisationID int64) (err error) {

	if coreValue.Text == "" {
		err = apperrors.TextFieldBlank
	}
	if coreValue.Description == "" {
		err = apperrors.DescFieldBlank
	}
	if coreValue.ParentID != nil {
		if !validateParentCoreValue(ctx, storer, organisationID, *coreValue.ParentID) {
			err = apperrors.InvalidParentValue
		}
	}

	return
}

func VarsStringToInt(inp string, label string) (result int64, err error) {

	if len(inp) <= 0 {
		err = apperrors.InvalidOrgId
		return
	}
	result, err = strconv.ParseInt(inp, 10, 64)
	if err != nil {
		logger.WithField("err", err.Error()).Error(fmt.Scanf("Error while parsing %s from url", label))
		err = apperrors.InternalServerError
		return

	}

	return
}
