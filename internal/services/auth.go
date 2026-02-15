package services

import (
	"context"
	"sport-assistance/internal/handlers/requests"
	"sport-assistance/internal/handlers/responses"
	"sport-assistance/internal/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Register(ctx context.Context, req requests.CreateUserRequest) (responses.JWTResponse, error) {
	birthDate, err := time.Parse(time.DateOnly, req.BirthDate)
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

	userId, err := s.repository.CreateUser(context.Background(), user)
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
