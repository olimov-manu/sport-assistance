package models

import "github.com/golang-jwt/jwt/v5"

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

type CustomClaims struct {
	UserId     int64             `json:"user_id"`
	Email      string            `json:"email"`
	Privileges map[string]string `json:"privileges"`
	jwt.RegisteredClaims
}

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
