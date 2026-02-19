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
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Register(ctx context.Context, req requests.CreateUserRequest) (responses.JWTResponse, error) {
	birthDate, err := time.Parse(s.cfg.DatabaseConfig.DBDateFormat, req.BirthDate)
	if err != nil {
		return responses.JWTResponse{}, myerrors.NewParseErr(myerrors.ParsingDateErrorMessage, err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return responses.JWTResponse{}, err
	}

	user := models.User{
		Name:                 req.Name,
		Surname:              req.Surname,
		Gender:               req.Gender,
		BirthDate:            birthDate,
		HeightCm:             req.HeightCm,
		WeightKg:             req.WeightKg,
		SportActivityLevelID: req.SportActivityLevelID,
		TownID:               req.TownID,
		RoleID:               req.RoleID,
		PhoneNumber:          req.PhoneNumber,
		Email:                req.Email,
		Password:             string(hash),
		IsHaveInjury:         req.IsHaveInjury,
		InjuryDescription:    req.InjuryDescription,
		Photo:                req.Photo,
	}

	userId, err := s.repository.CreateUser(ctx, user)
	if err != nil {
		return responses.JWTResponse{}, err
	}

	access, refresh, err := s.CreateTokens(ctx, userId, req.Email)
	if err != nil {
		return responses.JWTResponse{}, err
	}

	return responses.JWTResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

// Переписать логику обработки логина на номер телефона, не через email
func (s *Service) Login(ctx context.Context, req requests.LoginRequest) (responses.JWTResponse, error) {
	email := strings.TrimSpace(req.Email)
	user, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return responses.JWTResponse{}, myerrors.NewRepositoryErr(myerrors.UserDoesNotExistErrorMessage, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		fmt.Printf("bcrypt err: %T %v, passLen=%d\n", err, err, len(req.Password))
		return responses.JWTResponse{}, errors.New("invalid email or password")
	}

	accessToken, refreshToken, err := s.CreateTokens(ctx, user.ID, user.Email)
	if err != nil {
		return responses.JWTResponse{}, err
	}

	return responses.JWTResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Logout(ctx context.Context, req requests.LogoutRequest) (responses.EmptyResponse, error) {
	if req.RefreshToken == "" {
		return responses.EmptyResponse{}, errors.New("refresh token is required")
	}

	claims := &models.CustomClaims{}
	parsedToken, err := jwt.ParseWithClaims(req.RefreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, myerrors.NewTokenErr(myerrors.InvalidTokenErrorMessage, errors.New("invalid signing method"))
		}
		return []byte(s.cfg.SecurityConfig.RefreshTokenSecret), nil
	})
	if err != nil || parsedToken == nil || !parsedToken.Valid || claims.Subject != commons.RefreshSubject {
		return responses.EmptyResponse{}, errors.New("invalid refresh token")
	}

	if claims.UserId != req.UserID {
		return responses.EmptyResponse{}, errors.New("refresh token does not belong to this user")
	}

	if err := s.repository.RevokeRefreshToken(ctx, req.RefreshToken); err != nil {
		return responses.EmptyResponse{}, err
	}

	key := fmt.Sprintf(s.cfg.SecurityConfig.AccessTokenRedisPrefix, claims.UserId)
	fmt.Println(key)
	if err := s.redisClient.Del(ctx, key).Err(); err != nil {
		return responses.EmptyResponse{}, err
	}

	return responses.EmptyResponse{}, nil
}
