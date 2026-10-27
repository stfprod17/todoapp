package core_http_middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	core_auth "github.com/stfprod17/todoapp/internal/core/auth"
)

func Auth(tokenManager *core_auth.TokenManager) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			tokenString, ok := strings.CutPrefix(authorization, "Bearer ")
			if !ok || tokenString == "" {
				unauthorized(w, "authorization token is required")
				return
			}

			claims, err := tokenManager.Parse(tokenString)
			if err != nil {
				unauthorized(w, "authorization token is invalid")
				return
			}

			ctx := core_auth.WithUserID(r.Context(), claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func unauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": message})
}
