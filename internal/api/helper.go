package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func getPaginationParams(req *http.Request) (int64, int64, error) {
	pageStr := req.URL.Query().Get("page")
	limitStr := req.URL.Query().Get("limit")

	fmt.Println("pagestr: ", pageStr, " limitstr: ", limitStr)
	var page int64
	var limit int64
	var err error

	if pageStr == "" {
		page = 1
		limit = 10
	} else {
		fmt.Println("Hello page limit")
		page, err = strconv.ParseInt(pageStr, 10, 64)
		if err != nil || page < 1 {
			return 0, 0, errors.New("invalid page parameter")
		}

		if limitStr == "" {
			fmt.Println("empty limit")
			limit = 10
		} else {
			limit, err = strconv.ParseInt(limitStr, 10, 64)
			if err != nil || limit < 1 {
				return 0, 0, errors.New("invalid limit parameter")
			}
		}
	}

	fmt.Println("page: ", page, " limit: ", limit)
	return page, limit, nil
}
