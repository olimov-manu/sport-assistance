package services

import (
	"context"
	"errors"
	"fmt"
	"sport-assistance/internal/handlers/responses"
	"sport-assistance/pkg/myerrors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	otpStubValue = "0000"
	otpTTL       = 5 * time.Minute
)

func (s *Service) SendOTP(ctx context.Context, identifier string) (responses.SendOTPResponse, error) {
	normalizedIdentifier := strings.TrimSpace(strings.ToLower(identifier))
	if normalizedIdentifier == "" {
		return responses.SendOTPResponse{}, myerrors.NewValidationError("identifier is required", errors.New("empty identifier"))
	}

	key := fmt.Sprintf(s.cfg.SecurityConfig.OtpRedisPrefix, normalizedIdentifier)
	if err := s.redisClient.Set(ctx, key, otpStubValue, otpTTL).Err(); err != nil {
		return responses.SendOTPResponse{}, myerrors.NewTokenErr("failed to save otp in redis", err)
	}

	return responses.SendOTPResponse{
		OTPSent: true,
		Message: "OTP sent",
	}, nil
}

func (s *Service) ConfirmOTP(ctx context.Context, identifier, otp string) (responses.ConfirmOTPResponse, error) {
	normalizedIdentifier := strings.TrimSpace(strings.ToLower(identifier))
	normalizedOTP := strings.TrimSpace(otp)

	if normalizedIdentifier == "" {
		return responses.ConfirmOTPResponse{}, myerrors.NewValidationError("identifier is required", errors.New("empty identifier"))
	}
	if normalizedOTP == "" {
		return responses.ConfirmOTPResponse{}, myerrors.NewValidationError("otp is required", errors.New("empty otp"))
	}

	key := fmt.Sprintf(s.cfg.SecurityConfig.OtpRedisPrefix, normalizedIdentifier)
	savedOTP, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		return responses.ConfirmOTPResponse{}, myerrors.NewValidationError("otp is invalid or expired", err)
	}
	if savedOTP != normalizedOTP {
		return responses.ConfirmOTPResponse{}, myerrors.NewValidationError("otp does not match", errors.New("mismatch otp"))
	}

	if err = s.redisClient.Del(ctx, key).Err(); err != nil {
		return responses.ConfirmOTPResponse{}, myerrors.NewTokenErr("failed to delete otp from redis", err)
	}

	user, err := s.repository.GetUserByPhone(ctx, normalizedIdentifier)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return responses.ConfirmOTPResponse{
				OTPConfirmed: true,
				IsRegistered: false,
				Message:      "OTP confirmed, user is not registered",
			}, nil
		}

		return responses.ConfirmOTPResponse{}, myerrors.NewRepositoryErr("failed to fetch user by phone", err)
	}

	accessToken, refreshToken, err := s.CreateTokens(ctx, user.ID, user.Email)
	if err != nil {
		return responses.ConfirmOTPResponse{}, err
	}

	return responses.ConfirmOTPResponse{
		OTPConfirmed: true,
		IsRegistered: true,
		Message:      "OTP confirmed",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
