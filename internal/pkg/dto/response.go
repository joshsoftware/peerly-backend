package dto

import (
	"encoding/json"
	"net/http"

	logger "github.com/sirupsen/logrus"
)

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Error interface{} `json:"error"`
}

type MessageObject struct {
	Message string `json:"message"`
}

type ErrorObject struct {
	Code string `json:"code"`
	MessageObject
	Fields map[string]string `json:"fields"`
}

func Repsonse(rw http.ResponseWriter, status int, responseBody interface{}) {
	respBytes, err := json.Marshal(responseBody)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while marshaling core values data")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(respBytes)
}
