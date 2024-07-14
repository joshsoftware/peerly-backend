package middleware

import (
	"context"
	"net/http"
	"slices"

	"github.com/dgrijalva/jwt-go"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/config"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/sirupsen/logrus"
)

func JwtAuthMiddleware(next http.Handler, roles []string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// next.ServeHTTP(rw, req)
		// return 
		jwtKey := config.JWTKey()
		// authToken := req.Header.Get(constants.AuthorizationHeader)
		// if authToken == "" {
		// 	logger.Error("Empty auth token")
		// 	err := apperrors.InvalidAuthToken
		// 	dto.ErrorRepsonse(rw, err)
		// 	return
		// }

		authToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MSwiUm9sZSI6InVzZXIiLCJleHAiOjE3MjMxODgwODl9.K3FL_2BhYiIyf8tt_CZjTTTIbu_UVbNahc5hqby2dM8"
		claims := &dto.Claims{}

		tkn, err := jwt.ParseWithClaims(authToken, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error in parse with claims function")
			err = apperrors.InvalidAuthToken
			dto.ErrorRepsonse(rw, err)
			return
		}

		if !tkn.Valid {
			logger.Error("Invalid token")
			err = apperrors.InvalidAuthToken
			dto.ErrorRepsonse(rw, err)
			return
		}

		var Id int64 = claims.Id
		var Role string = claims.Role

		if !slices.Contains(roles, Role) {
			err := apperrors.RoleUnathorized
			dto.ErrorRepsonse(rw, err)
			return
		}

		// set id and role to context
		ctx := context.WithValue(req.Context(), constants.UserId, Id)
		ctx = context.WithValue(ctx, constants.Role, Role)
		req = req.WithContext(ctx)

		next.ServeHTTP(rw, req)

	})
}

func RecoverMiddleware(next http.Handler, roles []string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			rvr := recover()
			if rvr != nil {
				logger.Error("err ", rvr)
				dto.ErrorRepsonse(rw, apperrors.InternalServer)
			}
		}()

		next.ServeHTTP(rw, req)
	})
}
