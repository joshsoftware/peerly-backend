package api

import (
	"errors"
	"net/http"
	"strconv"
)

func getPaginationParams(req *http.Request) (int64, int64, error) {
	
	pageStr := req.URL.Query().Get("page")
	limitStr := req.URL.Query().Get("limit")
	var page int64
	var limit int64
	var err error

	if pageStr == "" {
		page = 1
		limit = 10
	} else {

		page, err = strconv.ParseInt(pageStr, 10, 64)
		if err != nil || page < 1 {
			return 0, 0, errors.New("invalid page parameter")
		}

		if limitStr == "" {
			limit = 10
		} else {
			limit, err = strconv.ParseInt(limitStr, 10, 64)
			if err != nil || limit < 1 {
				return 0, 0, errors.New("invalid limit parameter")
			}
		}
	}

	//TODO : max limit and min limit
	
	return page, limit, nil
}
