package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func newClaims(userID int64, ttl time.Duration) *Claims {
	now := time.Now()
	return &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)), // время истечения
			IssuedAt:  jwt.NewNumericDate(now),          // время выпуска
			// Можно добавить Issuer, Subject и т.п. если нужно
		},
	}
}

func GenerateToken(userID int64, secretKey string) (string, error) {
	claims := newClaims(userID, 24*time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
