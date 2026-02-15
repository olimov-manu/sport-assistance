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
	CreateUser(ctx context.Context, user models.User) (int64, error)
	CreateRefreshToken(ctx context.Context, userID int64, refreshToken string, expiresAt time.Time) error
	UserExistsByEmail(ctx context.Context, email string) (bool, error)
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
