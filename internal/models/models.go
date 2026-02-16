package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

type CustomClaims struct {
	UserId     uint64            `json:"user_id"`
	Email      string            `json:"email"`
	Privileges map[string]string `json:"privileges"`
	jwt.RegisteredClaims
}

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	ID        uint64     `json:"id"`
	UserID    uint64     `json:"user_id"`
	Token     string     `json:"token"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at"`
}
