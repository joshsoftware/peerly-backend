package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/config"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	log "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
	"github.com/sirupsen/logrus"
)


func loggerHandler(rw http.ResponseWriter, req *http.Request) {

	log.Debug(req.Context(),"loggerHandler: req: ",req)
	var changeLogRequest dto.ChangeLogLevelRequest
	err := json.NewDecoder(req.Body).Decode(&changeLogRequest)
	if err != nil {
		log.Errorf(req.Context(),"Error while decoding request data : %v", err)
		err = apperrors.JSONParsingErrorReq
		dto.ErrorRepsonse(rw, err)
		return
	}

	if config.DeveloperKey() != changeLogRequest.DeveloperKey {
		dto.ErrorRepsonse(rw,apperrors.UnauthorizedDeveloper)
		return 
	}

	log.Info(req.Context(), "loggerHandler")
	if changeLogRequest.LogLevel == "DebugLevel" {
		log.Logger.SetLevel(logrus.DebugLevel)
	}else if changeLogRequest.LogLevel == "InfoLevel" {
		log.Logger.SetLevel(logrus.InfoLevel)
	}else {
		dto.ErrorRepsonse(rw, apperrors.InvalidLoggerLevel)
		return
	}
	dto.SuccessRepsonse(rw, http.StatusOK, "Success", fmt.Sprintf("log level changed to %s", changeLogRequest.LogLevel))
}
