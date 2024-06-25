package api

import (
	"net/http"

	"github.com/joshsoftware/peerly-backend/internal/api/validation"
	user "github.com/joshsoftware/peerly-backend/internal/app/users"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

func loginUser(userSvc user.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		authToken := req.Header.Get("Authorization")

		validateResp, err := userSvc.ValidatePeerly(req.Context(), authToken)
		if err != nil {
			apperrors.ErrorResp(rw, err)
			return
		}

		reqData := dto.GetIntranetUserDataReq{
			Token:  validateResp.PeerlyToken,
			UserId: validateResp.UserId,
		}

		user, err := userSvc.GetIntranetUserData(req.Context(), reqData)
		if err != nil {
			apperrors.ErrorResp(rw, err)
			return
		}

		err = validation.GetIntranetUserDataValidation(user)
		if err != nil {
			apperrors.ErrorResp(rw, err)
			return
		}

		resp, err := userSvc.LoginUser(req.Context(), user)
		if err != nil {
			apperrors.ErrorResp(rw, err)
			return
		}

		dto.Repsonse(rw, http.StatusOK, dto.SuccessResponse{Data: resp})

	}
}
