package utils

import (
	"net/http"
	"strconv"
	"time"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	logger "github.com/sirupsen/logrus"
)

func VarsStringToInt(inp string, label string) (result int64, err error) {

	if len(inp) <= 0 {
		err = apperrors.InvalidId
		return
	}
	result, err = strconv.ParseInt(inp, 10, 64)
	if err != nil {
		logger.Errorf("error while parsing %s from url, err: %s", label, err.Error())
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

func GetPaginationParams(req *http.Request) (page int16, limit int16) {

	pageStr := req.URL.Query().Get("page")
	pageSizeStr := req.URL.Query().Get("page_size")

	page = constants.DefaultPageNumber
	if pageStr != "" {
		pageNumber, err := strconv.ParseInt(pageStr, 10, 32)
		if err != nil {
			logger.Errorf("err: %v", err)
		} else if pageNumber > 0 {
			page = int16(pageNumber)
		}
	}

	limit = constants.DefaultPageSize
	if pageSizeStr != "" {
		pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
		if err != nil {
			logger.Errorf("err: %v", err)
		} else if pageSize > constants.MaxPageSize {
			pageSize = constants.MaxPageSize
		}
		limit = int16(pageSize)
	}

	return page, limit
}
func GetSelfParam(req *http.Request) bool {
	paramStr := req.URL.Query().Get("self")
	if paramStr == "" {
		return false
	}

	boolValue, err := strconv.ParseBool(paramStr)
	if err != nil {
		logger.Errorf("err: %v", err)
		return false
	}

	return boolValue
}
