package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

const (
	// JWTTokenHeaderKey is the key for JWT token in the request header
	JWTTokenHeaderKey = "Authorization"
)

type contextKey string

const (
	userIDKey contextKey = "userID"
)

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(JWTTokenHeaderKey)
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(s.cfg.SecretKey), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok || userID == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Добавляем userID в контекст запроса
		ctx := r.Context()
		ctx = context.WithValue(ctx, userIDKey, int64(userID))
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
