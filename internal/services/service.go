package services

import (
	"context"
	"log/slog"
	"sport-assistance/internal/models"
	"sport-assistance/pkg/configs"
	"time"

	"github.com/redis/go-redis/v9"
)

type IRepository interface {
	CreateUser(ctx context.Context, user models.User) (uint64, error)
	GetUserByID(ctx context.Context, userID uint64) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	UserExistsByEmail(ctx context.Context, email string) (bool, error)
	RotateRefreshToken(ctx context.Context, userID uint64, oldRefreshToken, newRefreshToken string, newExpiresAt time.Time) error
	CreateRefreshToken(ctx context.Context, userID uint64, refreshToken string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, refreshToken string) (models.RefreshTokenResponse, error)
}

type Service struct {
	repository  IRepository
	logger      *slog.Logger
	cfg         *configs.Config
	redisClient *redis.Client
}

func NewService(repo IRepository, log *slog.Logger, cfg *configs.Config, redisClient *redis.Client) *Service {
	return &Service{
		repository:  repo,
		logger:      log,
		cfg:         cfg,
		redisClient: redisClient,
	}
}
