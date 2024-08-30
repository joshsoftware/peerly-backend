package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/config"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/sirupsen/logrus"
)

func JwtAuthMiddleware(next http.Handler, role int) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		jwtKey := config.JWTKey()
		authToken := req.Header.Get(constants.AuthorizationHeader)
		if authToken == "" {
			logger.Error("Empty auth token")
			err := apperrors.InvalidAuthToken
			dto.ErrorRepsonse(rw, err)
			return
		}

		authToken = strings.TrimPrefix(authToken, "Bearer ")
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

		Id := claims.Id
		Role := claims.Role

		if Role > role {
			err := apperrors.RoleUnathorized
			dto.ErrorRepsonse(rw, err)
			return
		}

		// set id and role to context
		fmt.Println("setting id: ", Id)
		ctx := context.WithValue(req.Context(), constants.UserId, Id)
		ctx = context.WithValue(ctx, constants.Role, Role)
		req = req.WithContext(ctx)

		next.ServeHTTP(rw, req)

	})
}
