package api

import (
	"fmt"
	"net/http"
	"strconv"

	logger "github.com/sirupsen/logrus"
)

func getPaginationParams(req *http.Request) (page int16, limit int16) {

	pageStr := req.URL.Query().Get("page")
	limitStr := req.URL.Query().Get("page_size")

	if pageStr == "" {
		page = 1
	} else {
		pageInt64, err := strconv.ParseInt(pageStr, 10, 32)
		if err != nil {
			logger.Error(fmt.Sprintf("err: %v",err))
			page = 1 
		}

		if pageInt64 < 1 {
			pageInt64 = 1
		}
		page = int16(pageInt64)
	}

	if limitStr == "" {
		limit = 10
	} else {
		limitInt64, err := strconv.ParseInt(limitStr, 10, 16)
		if err != nil  {
			logger.Error(fmt.Sprintf("err: %v",err))
			limit = 10
		}
		if limitInt64 < 1 {
			limitInt64 = 10
		}else if limitInt64 > 1000 {
			limitInt64 = 1000
		}
		limit = int16(limitInt64)
	}

	return page, limit
}
func getSelfParam(req *http.Request) (bool) {
	paramStr := req.URL.Query().Get("self")
	if paramStr == "" {
		return false
	}
	
	boolValue, err := strconv.ParseBool(paramStr)
	if err != nil {
		return false
	}
	
	return boolValue
}
