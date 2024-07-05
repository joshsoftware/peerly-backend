package dto

import (
	"encoding/json"
	"net/http"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	logger "github.com/sirupsen/logrus"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Status  int         `json:"status_code"`
	Data    interface{} `json:"data"`
}

func SuccessRepsonse(rw http.ResponseWriter, status int, message string, data interface{}) {

	var resp Response
	resp.Success = true
	resp.Status = status
	resp.Message = message
	resp.Data = data

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while marshaling core values data")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(respBytes)
}

func ErrorRepsonse(rw http.ResponseWriter, err error) {

	var resp Response
	resp.Success = false
	resp.Status = apperrors.GetHTTPStatusCode(err)
	resp.Message = err.Error()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while marshaling core values data")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(resp.Status)
	rw.Write(respBytes)
}
