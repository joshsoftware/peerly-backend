package api

import (
	"net/http"
	"strconv"

	logger "github.com/sirupsen/logrus"
)

func getPaginationParams(req *http.Request) (page int16, limit int16) {

	pageStr := req.URL.Query().Get("page")
	limitStr := req.URL.Query().Get("page_size")

	// if pageStr == "" {
	// 	page = 1
	// } else {
	// 	pageInt64, err := strconv.ParseInt(pageStr, 10, 32)
	// 	if err != nil {
	// 		logger.Errorf("err: %v",err)
	// 		page = 1
	// 	}

	// 	if pageInt64 < 1 {
	// 		pageInt64 = 1
	// 	}
	// 	page = int16(pageInt64)
	// }

	page = 1
	if pageStr != "" {
		pageInt64, err := strconv.ParseInt(pageStr, 10, 32)
		if err != nil {
			logger.Errorf("err: %v", err)
		}else if pageInt64 > 0 {
			page = int16(pageInt64)
		}
	}

	limit = 10

	if limitStr != "" {
		limitInt64, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil {
			logger.Errorf("err: %v", err)
		}else if limitInt64 > 1000 {
			limitInt64 = 1000
		}
		limit = int16(limitInt64)
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
