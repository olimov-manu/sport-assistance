package repositories

import (
	"context"
	"time"
)

func (r *Repository) CreateRefreshToken(ctx context.Context, userID int64, refreshToken string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (
			user_id,
			token,
			expires_at
		)
		VALUES ($1, $2, $3)
	`

	_, err := r.postgres.Exec(
		ctx,
		query,
		userID,
		refreshToken,
		expiresAt,
	)

	return err
}
