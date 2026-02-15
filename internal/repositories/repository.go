package repositories

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	postgres *pgxpool.Pool
	logger   *slog.Logger
}

func NewRepository(postgresConn *pgxpool.Pool, log *slog.Logger) *Repository {
	return &Repository{
		postgres: postgresConn,
		logger:   log,
	}
}
