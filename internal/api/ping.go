package api

import (
	"net/http"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	log "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
)

func pingHandler(rw http.ResponseWriter, req *http.Request) {
	log.Debug(req.Context(),"debug ping")
	log.Info(req.Context(),"ping")
	dto.SuccessRepsonse(rw, http.StatusOK, "Success", "pong")
}
