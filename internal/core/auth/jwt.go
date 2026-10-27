package core_auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stfprod17/todoapp/internal/core/domain"
)

type TokenManager struct {
	secret []byte
	ttl    time.Duration
}

type Claims struct {
	UserID int
	Email  string
	jwt.RegisteredClaims
}

func NewTokenManager(secret string) (*TokenManager, error) {
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is empty")
	}

	return &TokenManager{secret: []byte(secret), ttl: 24 * time.Hour}, nil
}

func (m *TokenManager) Parse(tokenString string) (Claims, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return Claims{}, fmt.Errorf("parse token: %w", err)
	}
	if !token.Valid {
		return Claims{}, fmt.Errorf("token is invalid")
	}
	return claims, nil
}

func (m *TokenManager) Generate(user domain.User) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}
