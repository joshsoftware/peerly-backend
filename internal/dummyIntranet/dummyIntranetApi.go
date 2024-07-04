package intranet

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type User struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	ProfileImgUrl string `json:"profile_image_url"`
	Designation   string `json:"designation"`
	Grade         string `json:"grade"`
}

type ValidateResp struct {
	PeerlyToken string `json:"peerly_token"`
	UserId      int    `json:"user_id"`
}

func ValidatePeerly() http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		authToken := req.Header.Get("Authorization")
		peerlyCode := req.Header.Get("PeerlyCode")
		if authToken != "peerly" || peerlyCode != "peerly" {
			err := apperrors.InvalidAuthToken
			dto.ErrorRepsonse(rw,err,nil)
			return
		}
		resp := ValidateResp{
			PeerlyToken: "peerly",
			UserId:      3,
		}
		respBody, err := json.Marshal(resp)
		if err != nil {
			err := apperrors.JSONParsingErrorResp
			dto.ErrorRepsonse(rw,err,nil)
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write(respBody)
	})
}

func IntranetGetUserApi() http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		authToken := req.Header.Get("Authorization")
		vars := mux.Vars(req)
		if vars["user_id"] == "" {
			err := apperrors.InternalServerError
			dto.ErrorRepsonse(rw,err,nil)
			return
		}
		fmt.Println("response for user id: ", vars["user_id"])
		if authToken != "peerly" {
			logger.WithField("err", "err").Error("Error in authtoken! Authtoken = ", authToken)
			err := apperrors.InvalidAuthToken
			dto.ErrorRepsonse(rw,err,nil)
			return
		}
		user := User{
			FirstName:     "Sharyu",
			LastName:      "Marwadi",
			Email:         "sharyu2@joshsoftware.com",
			ProfileImgUrl: "imgurl",
			Designation:   "Intern",
			Grade:         "J12",
		}
		resp, err := json.Marshal(user)
		if err != nil {
			err := apperrors.JSONParsingErrorResp
			dto.ErrorRepsonse(rw,err,nil)
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write(resp)
	})
}
