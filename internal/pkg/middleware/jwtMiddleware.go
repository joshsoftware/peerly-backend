package middleware

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/config"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/sirupsen/logrus"
)

func JwtAuthMiddleware(next http.Handler, roles []string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		jwtKey := config.JWTKey()
		logger.Infof("sp: jwtkey: %v",jwtKey)
		authToken := req.Header.Get(constants.AuthorizationHeader)
		logger.Infof("sp: authtoken: %v",authToken)
		if authToken == "" {
			logger.Error("Empty auth token")
			err := apperrors.InvalidAuthToken
			dto.ErrorRepsonse(rw, err)
			return
		}

		authToken = strings.TrimPrefix(authToken, "Bearer ")
		logger.Infof("sp: trimAuthToken: %v",authToken)
		claims := &dto.Claims{}

		tkn, err := jwt.ParseWithClaims(authToken, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		logger.Infof("sp: Token: %v",tkn)
		logger.Infof("sp: Token err: %v",err)
		
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

		Id := claims.Id
		Role := claims.Role

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
