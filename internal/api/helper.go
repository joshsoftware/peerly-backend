package api

import (
	"net/http"
	"strconv"

	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	logger "github.com/sirupsen/logrus"
)

func getPaginationParams(req *http.Request) (page int16, limit int16) {

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
func getSelfParam(req *http.Request) bool {
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
