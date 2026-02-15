package services

import (
	"context"
	"errors"
	"fmt"
	"sport-assistance/internal/models"
	"sport-assistance/pkg/commons"
	"sport-assistance/pkg/myerrors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTIssuer struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
}

func (s *Service) CreateTokens(ctx context.Context, userID int64, email string) (string, string, error) {
	now := time.Now()

	accessToken, err := s.createAccessToken(userID, email, now)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.createRefreshToken(ctx, userID, now)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *Service) createAccessToken(userID int64, email string, now time.Time) (string, error) {
	claims := s.buildClaims(userID, email, now, s.cfg.SecurityConfig.AccessTokenTTL, commons.AccessSubject)

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.SecurityConfig.AccessTokenSecret))
	if err != nil {
		return "", myerrors.NewTokenErr(myerrors.AccessTokenCreateErrorMessage, err)
	}

	return token, nil
}

func (s *Service) createRefreshToken(ctx context.Context, userID int64, now time.Time) (string, error) {
	expiresAt := now.Add(s.cfg.SecurityConfig.RefreshTokenTTL)
	claims := s.buildClaims(userID, "", now, s.cfg.SecurityConfig.RefreshTokenTTL, commons.RefreshSubject)

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.SecurityConfig.RefreshTokenSecret))
	if err != nil {
		return "", myerrors.NewTokenErr(myerrors.RefreshTokenCreateErrorMessage, err)
	}

	if err = s.repository.CreateRefreshToken(ctx, userID, token, expiresAt); err != nil {
		return "", myerrors.NewTokenErr(myerrors.RefreshTokenCreateInDBErrorMessage, err)
	}

	if err = s.saveToRedis(ctx, userID, token, expiresAt); err != nil {
		return "", myerrors.NewTokenErr(myerrors.RefreshTokenCreateInRedisErrorMessage, err)
	}

	return token, nil
}

func (s *Service) buildClaims(userID int64, email string, now time.Time, ttl time.Duration, subject string) models.CustomClaims {
	privileges := map[string]string{}

	return models.CustomClaims{
		UserId:     userID,
		Email:      email,
		Privileges: privileges,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        fmt.Sprintf("uid:%d", userID),
		},
	}
}

func (s *Service) saveToRedis(ctx context.Context, userID int64, token string, expiresAt time.Time) error {
	key := strconv.FormatInt(userID, 10)
	ttl := time.Until(expiresAt)
	return s.redisClient.Set(ctx, key, token, ttl).Err()
}

func ValidateJWT(tokenString string, secret string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if uid, ok := claims["user_id"].(float64); ok {
			return int(uid), nil
		}
	}

	return 0, errors.New("invalid token")
}
