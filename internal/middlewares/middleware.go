package middlewares

import (
	"log/slog"
	"sport-assistance/internal/services"
	"sport-assistance/pkg/configs"

	"github.com/redis/go-redis/v9"
)

type Middleware struct {
	repo        services.IRepository
	cfg         configs.SecurityConfig
	logger      *slog.Logger
	redisClient *redis.Client
}

func NewMiddleware(repo services.IRepository, cfg configs.SecurityConfig, log *slog.Logger, redisClient *redis.Client) *Middleware {
	return &Middleware{
		cfg:         cfg,
		repo:        repo,
		logger:      log,
		redisClient: redisClient,
	}
}
