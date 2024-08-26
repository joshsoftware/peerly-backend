package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/peerly-backend/internal/app/reward"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	log "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
)


func giveRewardHandler(rewardSvc reward.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		apprId, err := strconv.Atoi(vars["id"])
		if err != nil {
			dto.ErrorRepsonse(rw, apperrors.BadRequest)
			return
		}

		log.Debug(req.Context(),"giveRewardHandler: request: ",req)

		var reward dto.Reward
		err = json.NewDecoder(req.Body).Decode(&reward)
		if err != nil {
			log.Error(req.Context(),"Error decoding request data:", err.Error())
			dto.ErrorRepsonse(rw, apperrors.JSONParsingErrorReq)
			return
		}
		

		if reward.Point <1 || reward.Point >5 {
			log.Error(req.Context(),"Invalid reward point")
			dto.ErrorRepsonse(rw,apperrors.InvalidRewardPoint)
			return 
		}
		reward.AppreciationId = int64(apprId)
		resp, err := rewardSvc.GiveReward(req.Context(),reward)
		if err != nil {
			log.Error(req.Context(),"resp err: ",err)
			dto.ErrorRepsonse(rw, err)
			return
		}

		log.Debug(req.Context(),"giveRewardHandler: resp: ",resp)
		
		dto.SuccessRepsonse(rw, http.StatusCreated, "Reward given successfully", resp)
	})
}
