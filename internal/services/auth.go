package services

import (
	"context"
	"errors"
	"fmt"
	"sport-assistance/internal/handlers/requests"
	"sport-assistance/internal/handlers/responses"
	"sport-assistance/internal/models"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Register(ctx context.Context, req requests.CreateUserRequest) (responses.JWTResponse, error) {
	birthDate, err := time.Parse(s.cfg.DatabaseConfig.DBDateFormat, req.BirthDate)
	if err != nil {
		return responses.JWTResponse{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return responses.JWTResponse{}, err
	}

	user := models.User{
		Name:      req.Name,
		Surname:   req.Surname,
		Gender:    req.Gender,
		BirthDate: birthDate,

		HeightCm:             req.HeightCm,
		WeightKg:             req.WeightKg,
		SportActivityLevelID: req.SportActivityLevelID,
		TownID:               req.TownID,

		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Password:    string(hash),

		IsHaveInjury:      req.IsHaveInjury,
		InjuryDescription: req.InjuryDescription,
		Photo:             req.Photo,
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

func (s *Service) Login(ctx context.Context, req requests.LoginRequest) (responses.JWTResponse, error) {
	email := strings.TrimSpace(req.Email)
	user, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return responses.JWTResponse{}, errors.New("this user does not exist")
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
