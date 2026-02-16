package services

import (
	"context"
	"errors"
	"fmt"
	"sport-assistance/internal/handlers/requests"
	"sport-assistance/internal/handlers/responses"
	"sport-assistance/internal/models"
	"sport-assistance/pkg/commons"
	"sport-assistance/pkg/myerrors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTIssuer struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
}

func (s *Service) CreateTokens(ctx context.Context, userID uint64, email string) (string, string, error) {
	now := time.Now()
	expiresAt := now.Add(s.cfg.SecurityConfig.RefreshTokenTTL)

	accessToken, err := s.createAccessToken(ctx, userID, email, now)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.createRefreshToken(userID, now)
	if err != nil {
		return "", "", err
	}

	if err = s.repository.CreateRefreshToken(ctx, userID, refreshToken, expiresAt); err != nil {
		return "", "", myerrors.NewTokenErr(myerrors.RefreshTokenCreateInDBErrorMessage, err)
	}

	return accessToken, refreshToken, nil
}

func (s *Service) RefreshTokens(ctx context.Context, req requests.RefreshTokensRequest) (responses.JWTResponse, error) {
	oldRefreshToken := req.RefreshToken
	claims := &models.CustomClaims{}
	now := time.Now()
	expiresAtRefreshToken := now.Add(s.cfg.SecurityConfig.RefreshTokenTTL)

	token, err := jwt.ParseWithClaims(req.RefreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.SecurityConfig.RefreshTokenSecret), nil
	})
	if err != nil || !token.Valid || claims.Subject != commons.RefreshSubject {
		return responses.JWTResponse{}, errors.New("invalid refresh token")
	}

	refreshToken, err := s.repository.GetRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return responses.JWTResponse{}, err
	}

	if refreshToken.RevokedAt != nil {
		return responses.JWTResponse{}, errors.New("refresh token is revoked")
	}

	if s.IsTokenExpired(refreshToken.ExpiresAt) {
		return responses.JWTResponse{}, errors.New("refresh token is expired")
	}

	user, err := s.repository.GetUserByID(ctx, refreshToken.UserID)
	if err != nil {
		return responses.JWTResponse{}, err
	}

	newRefreshToken, err := s.createRefreshToken(user.ID, now)
	if err != nil {
		return responses.JWTResponse{}, err
	}

	err = s.repository.RotateRefreshToken(ctx, user.ID, oldRefreshToken, newRefreshToken, expiresAtRefreshToken)
	if err != nil {
		return responses.JWTResponse{}, err
	}

	accessToken, err := s.createAccessToken(ctx, refreshToken.UserID, user.Email, now)
	if err != nil {
		return responses.JWTResponse{}, err
	}

	return responses.JWTResponse{AccessToken: accessToken, RefreshToken: newRefreshToken}, nil
}

func (s *Service) IsTokenExpired(expiresAt time.Time) bool {
	return time.Now().After(expiresAt)
}

func (s *Service) createAccessToken(ctx context.Context, userID uint64, email string, now time.Time) (string, error) {
	claims := s.buildClaims(userID, email, now, s.cfg.SecurityConfig.AccessTokenTTL, commons.AccessSubject)

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.SecurityConfig.AccessTokenSecret))
	if err != nil {
		return "", myerrors.NewTokenErr(myerrors.AccessTokenCreateErrorMessage, err)
	}

	key := fmt.Sprintf(s.cfg.SecurityConfig.AccessTokenRedisPrefix, userID)
	if err = s.saveToRedis(ctx, key, token, claims.ExpiresAt.Time); err != nil {
		return "", myerrors.NewTokenErr(myerrors.RefreshTokenCreateInRedisErrorMessage, err)
	}

	return token, nil
}

func (s *Service) createRefreshToken(userID uint64, now time.Time) (string, error) {
	claims := s.buildClaims(userID, "", now, s.cfg.SecurityConfig.RefreshTokenTTL, commons.RefreshSubject)

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.SecurityConfig.RefreshTokenSecret))
	if err != nil {
		return "", myerrors.NewTokenErr(myerrors.RefreshTokenCreateErrorMessage, err)
	}

	return token, nil
}

func (s *Service) buildClaims(userID uint64, email string, now time.Time, ttl time.Duration, subject string) models.CustomClaims {
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
			ID:        uuid.NewString(),
		},
	}
}

func (s *Service) saveToRedis(ctx context.Context, key, token string, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	return s.redisClient.Set(ctx, key, token, ttl).Err()
}
