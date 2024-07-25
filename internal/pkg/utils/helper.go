package utils

import (
	"fmt"
	"strconv"
	"time"

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

// GetQuarter returns financial quarter
func GetQuarter() int8 {
	month := int(time.Now().Month())
	if month >= 1 && month <= 3 {
		return 4
	} else if month >= 4 && month <= 6 {
		return 1
	} else if month >= 7 && month <= 9 {
		return 2
	} else if month >= 10 && month <= 12 {
		return 3
	}
	return -1
}