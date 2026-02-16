package repositories

import (
	"context"
	"errors"
	"sport-assistance/internal/models"
	"time"
)

func (r *Repository) CreateRefreshToken(ctx context.Context, userID uint64, refreshToken string, expiresAt time.Time) error {
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

func (r *Repository) GetRefreshToken(ctx context.Context, refreshToken string) (models.RefreshTokenResponse, error) {
	const q = `
        SELECT id, user_id, expires_at, revoked_at
        FROM refresh_tokens
        WHERE token = $1
        LIMIT 1
    `
	var refreshTokenResponse models.RefreshTokenResponse

	err := r.postgres.QueryRow(ctx, q, refreshToken).Scan(
		&refreshTokenResponse.ID,
		&refreshTokenResponse.UserID,
		&refreshTokenResponse.ExpiresAt,
		&refreshTokenResponse.RevokedAt,
	)

	if err != nil {
		return models.RefreshTokenResponse{}, err
	}
	return refreshTokenResponse, nil
}

func (r *Repository) RotateRefreshToken(
	ctx context.Context,
	userID uint64,
	oldRefreshToken string,
	newRefreshToken string,
	newExpiresAt time.Time,
) error {
	tx, err := r.postgres.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Атомарно: отзываем старый (если он валиден и не отозван) и только тогда вставляем новый.
	ct, err := tx.Exec(ctx, `
		WITH revoked AS (
			UPDATE refresh_tokens
			SET revoked_at = now()
			WHERE user_id = $1
			  AND token = $2
			  AND revoked_at IS NULL
			  AND expires_at > now()
			RETURNING 1
		)
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		SELECT $1, $3, $4
		WHERE EXISTS (SELECT 1 FROM revoked)
	`, userID, oldRefreshToken, newRefreshToken, newExpiresAt)
	if err != nil {
		return err
	}

	// Если старый токен не был отозван (не найден/просрочен/уже использован),
	// то вставка нового не произойдёт, и RowsAffected() будет 0.
	if ct.RowsAffected() == 0 {
		return errors.New("refresh token invalid, expired, or already used")
	}

	return tx.Commit(ctx)
}
