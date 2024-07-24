package utils

import (
	"fmt"
	"strconv"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	logger "github.com/sirupsen/logrus"
)

func VarsStringToInt(inp string, label string) (result int64, err error) {

	if len(inp) <= 0 {
		err = apperrors.InvalidId
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
