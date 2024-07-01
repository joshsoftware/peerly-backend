package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const UserIDKey contextKey = "userID"
func JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		// Create a new context with the user ID
		ctx := context.WithValue(r.Context(), UserIDKey, int64(4))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
