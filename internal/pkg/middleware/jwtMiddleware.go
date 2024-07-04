package middleware

import (
	"context"
	"net/http"
	"slices"

	"github.com/dgrijalva/jwt-go"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/config"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/sirupsen/logrus"
)

func JwtAuthMiddleware(next http.Handler, roles []string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// next.ServeHTTP(rw, req)
		jwtKey := config.JWTKey()
		authToken := req.Header.Get("Authorization")
		if authToken == "" {
			logger.Error("Empty auth token")
			err := apperrors.InvalidAuthToken
			dto.ErrorRepsonse(rw, err, nil)
			return
		}

		claims := &dto.Claims{}

		tkn, err := jwt.ParseWithClaims(authToken, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error in parse with claims function")
			err = apperrors.InvalidAuthToken
			dto.ErrorRepsonse(rw, err, nil)
			return
		}

		if !tkn.Valid {
			logger.Error("Invalid token")
			err = apperrors.InvalidAuthToken
			dto.ErrorRepsonse(rw, err, nil)
			return
		}

		Id := claims.Id
		Role := claims.Role

		if !slices.Contains(roles, Role) {
			err := apperrors.RoleUnathorized
			dto.ErrorRepsonse(rw, err, nil)
			return
		}

		// set id and role to context
		ctx := context.WithValue(req.Context(), "userId", Id)
		ctx = context.WithValue(ctx, "role", Role)
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
				dto.ErrorRepsonse(rw, apperrors.InternalServer, nil)
			}
		}()

		next.ServeHTTP(rw, req)
	})
}
