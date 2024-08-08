package api

import (
	"net/http"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

func pingHandler(rw http.ResponseWriter, req *http.Request) {
	dto.SuccessRepsonse(rw, http.StatusOK, "Success", "testing")
}
