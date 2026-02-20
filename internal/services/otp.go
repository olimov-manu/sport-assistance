package services

import (
	"context"
	"errors"
	"fmt"
	"sport-assistance/pkg/myerrors"
	"strings"
	"time"
)

const (
	otpStubValue = "0000"
	otpTTL       = 5 * time.Minute
)

func (s *Service) SendOTP(ctx context.Context, identifier string) (string, error) {
	normalizedIdentifier := strings.TrimSpace(strings.ToLower(identifier))
	if normalizedIdentifier == "" {
		return "", myerrors.NewValidationError("identifier is required", errors.New("empty identifier"))
	}

	key := fmt.Sprintf(s.cfg.SecurityConfig.OtpRedisPrefix, normalizedIdentifier)
	if err := s.redisClient.Set(ctx, key, otpStubValue, otpTTL).Err(); err != nil {
		return "", myerrors.NewTokenErr("failed to save otp in redis", err)
	}

	return otpStubValue, nil
}

func (s *Service) ConfirmOTP(ctx context.Context, identifier, otp string) error {
	normalizedIdentifier := strings.TrimSpace(strings.ToLower(identifier))
	normalizedOTP := strings.TrimSpace(otp)

	if normalizedIdentifier == "" {
		return myerrors.NewValidationError("identifier is required", errors.New("empty identifier"))
	}
	if normalizedOTP == "" {
		return myerrors.NewValidationError("otp is required", errors.New("empty otp"))
	}

	key := fmt.Sprintf(s.cfg.SecurityConfig.OtpRedisPrefix, normalizedIdentifier)
	savedOTP, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		return myerrors.NewValidationError("otp is invalid or expired", err)
	}
	if savedOTP != normalizedOTP {
		return myerrors.NewValidationError("otp does not match", errors.New("mismatch otp"))
	}

	if err = s.redisClient.Del(ctx, key).Err(); err != nil {
		return myerrors.NewTokenErr("failed to delete otp from redis", err)
	}

	return nil
}
