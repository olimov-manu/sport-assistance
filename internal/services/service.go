package services

import (
	"context"
	"log/slog"
	"sport-assistance/internal/models"
	"sport-assistance/internal/services/dto"
	"sport-assistance/pkg/configs"
	"time"

	"github.com/redis/go-redis/v9"
)

type IRepository interface {
	// User
	CreateUser(ctx context.Context, user models.User) (uint64, error)
	GetUsers(ctx context.Context) ([]dto.UserDto, error)
	GetUserByID(ctx context.Context, userID uint64) (dto.UserDto, error)
	GetUserByEmail(ctx context.Context, email string) (dto.UserDto, error)
	UpdateUser(ctx context.Context, userID uint64, user models.User) error
	DeleteUser(ctx context.Context, userID uint64) error
	UserExistsByEmail(ctx context.Context, email string) (bool, error)

	// Permissions
	GetPermissionsByRoleId(ctx context.Context, roleId uint64) ([]string, error)

	// Jwt Tokens
	CreateRefreshToken(ctx context.Context, userID uint64, refreshToken string, expiresAt time.Time) error
	RotateRefreshToken(ctx context.Context, userID uint64, oldRefreshToken, newRefreshToken string, newExpiresAt time.Time) error
	GetRefreshToken(ctx context.Context, refreshToken string) (models.RefreshTokenResponse, error)
	RevokeRefreshToken(ctx context.Context, refreshToken string) error
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
