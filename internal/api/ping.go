package api

import (
	"net/http"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	log "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
)

func pingHandler(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log.Debug(ctx, "debug ping")
	log.Info(ctx, "ping")
	dto.SuccessRepsonse(rw, http.StatusOK, "Success", "pong")
}
