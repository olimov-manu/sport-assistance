package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sport-assistance/internal/models"
	"sport-assistance/pkg/myerrors"
	"time"

	"github.com/jackc/pgx/v5"
)

// CreateRefreshToken создаёт новый refresh token в БД
func (r *Repository) CreateRefreshToken(ctx context.Context, userID uint64, refreshToken string, expiresAt time.Time) error {
	const q = `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`

	if err := ctx.Err(); err != nil {
		return myerrors.NewRepositoryErr("контекст отменён перед выполнением запроса: ", err)
	}

	_, err := r.postgres.Exec(ctx, q, userID, refreshToken, expiresAt)
	if err != nil {
		return myerrors.NewRepositoryErr("не удалось создать refresh token: ", err)
	}

	return nil
}

// GetRefreshToken получает refresh token по его значению
// Возвращает ErrRefreshTokenNotFound если токен не существует
func (r *Repository) GetRefreshToken(ctx context.Context, refreshToken string) (models.RefreshTokenResponse, error) {
	const q = `
		SELECT id, user_id, expires_at, revoked_at
		FROM refresh_tokens
		WHERE token = $1
	`

	if err := ctx.Err(); err != nil {
		return models.RefreshTokenResponse{}, myerrors.NewRepositoryErr("контекст отменён перед выполнением запроса: ", err)
	}

	var res models.RefreshTokenResponse

	err := r.postgres.QueryRow(ctx, q, refreshToken).Scan(
		&res.ID,
		&res.UserID,
		&res.ExpiresAt,
		&res.RevokedAt,
	)
	if err != nil {
		// Токен не найден — это нормальный бизнес-кейс, не ошибка БД
		if errors.Is(err, pgx.ErrNoRows) {
			return models.RefreshTokenResponse{}, myerrors.ErrRefreshTokenNotFound
		}

		return models.RefreshTokenResponse{}, myerrors.NewRepositoryErr(
			"не удалось получить refresh token: ",
			err,
		)
	}

	return res, nil
}

// RotateRefreshToken выполняет ротацию refresh token в рамках транзакции
// Отзывает старый токен и создаёт новый
// Возвращает ErrRefreshTokenInvalid если старый токен невалидный/просрочен/уже отозван
func (r *Repository) RotateRefreshToken(
	ctx context.Context,
	userID uint64,
	oldRefreshToken string,
	newRefreshToken string,
	newExpiresAt time.Time,
) error {
	if err := ctx.Err(); err != nil {
		return myerrors.NewRepositoryErr("контекст отменён перед выполнением запроса: ", err)
	}

	tx, err := r.postgres.Begin(ctx)
	if err != nil {
		return myerrors.NewRepositoryErr("не удалось начать транзакцию: ", err)
	}

	// Гарантируем откат при выходе, если транзакция не была закоммичена
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			// Игнорируем ошибку отката, если транзакция уже закончена
			if !errors.Is(err, pgx.ErrTxClosed) {
				log.Printf("ошибка отката транзакции: %v\n", err)
			}
		}
	}()

	// Отзываем старый токен и создаём новый в одной CTE-операции
	ct, err := tx.Exec(ctx, `
		WITH revoked AS (
			UPDATE refresh_tokens
			SET revoked_at = now()
			WHERE user_id = $1
			  AND token = $2
			  AND revoked_at IS NULL
			  AND expires_at > now()
			RETURNING id
		)
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		SELECT $1, $3, $4
		WHERE EXISTS (SELECT 1 FROM revoked)
	`, userID, oldRefreshToken, newRefreshToken, newExpiresAt)

	if err != nil {
		return myerrors.NewRepositoryErr("не удалось выполнить операцию ротации токена: ", err)
	}

	// Бизнес-кейс: старый токен невалиден
	// (не найден, просрочен или уже отозван)
	if ct.RowsAffected() == 0 {
		return myerrors.ErrRefreshTokenInvalid
	}

	// Коммитим транзакцию
	if err := tx.Commit(ctx); err != nil {
		return myerrors.NewRepositoryErr("не удалось закоммитить транзакцию: ", err)
	}

	// После успешного коммита defer не выполнит откат, так как
	// Rollback на закрытой транзакции вернёт ErrTxClosed
	return nil
}

// RevokeRefreshToken отзывает refresh token по его ID
func (r *Repository) RevokeRefreshToken(ctx context.Context, tokenID uint64) error {
	const q = `
		UPDATE refresh_tokens
		SET revoked_at = now()
		WHERE id = $1 AND revoked_at IS NULL
	`

	if err := ctx.Err(); err != nil {
		return myerrors.NewRepositoryErr("контекст отменён перед выполнением запроса: ", err)
	}

	ct, err := r.postgres.Exec(ctx, q, tokenID)
	if err != nil {
		return myerrors.NewRepositoryErr("не удалось отозвать refresh token: ", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("токен не найден или уже отозван")
	}

	return nil
}

// RevokeAllUserRefreshTokens отзывает все refresh токены пользователя
func (r *Repository) RevokeAllUserRefreshTokens(ctx context.Context, userID uint64) error {
	const q = `
		UPDATE refresh_tokens
		SET revoked_at = now()
		WHERE user_id = $1 AND revoked_at IS NULL
	`

	if err := ctx.Err(); err != nil {
		return myerrors.NewRepositoryErr("контекст отменён перед выполнением запроса: ", err)
	}

	_, err := r.postgres.Exec(ctx, q, userID)
	if err != nil {
		return myerrors.NewRepositoryErr("не удалось отозвать все refresh токены пользователя: ", err)
	}

	return nil
}

// IsRefreshTokenValid проверяет валидность refresh token
// Токен валиден, если:
// - существует
// - не просрочен
// - не отозван
func (r *Repository) IsRefreshTokenValid(ctx context.Context, refreshToken string) (bool, error) {
	const q = `
		SELECT EXISTS(
			SELECT 1
			FROM refresh_tokens
			WHERE token = $1
			  AND revoked_at IS NULL
			  AND expires_at > now()
		)
	`

	if err := ctx.Err(); err != nil {
		return false, myerrors.NewRepositoryErr("контекст отменён перед выполнением запроса: ", err)
	}

	var exists bool
	err := r.postgres.QueryRow(ctx, q, refreshToken).Scan(&exists)
	if err != nil {
		return false, myerrors.NewRepositoryErr("не удалось проверить валидность токена: ", err)
	}

	return exists, nil
}

// CleanupExpiredTokens удаляет истёкшие и отозванные токены
// (обычно вызывается периодически по расписанию)
func (r *Repository) CleanupExpiredTokens(ctx context.Context, olderThan time.Duration) (int64, error) {
	const q = `
		DELETE FROM refresh_tokens
		WHERE (revoked_at IS NOT NULL AND revoked_at < now() - interval '1 second' * $1)
		   OR (expires_at < now() - interval '1 second' * $2)
	`

	if err := ctx.Err(); err != nil {
		return 0, myerrors.NewRepositoryErr("контекст отменён перед выполнением запроса: ", err)
	}

	ct, err := r.postgres.Exec(ctx, q, int64(olderThan.Seconds()), int64(olderThan.Seconds()))
	if err != nil {
		return 0, myerrors.NewRepositoryErr("не удалось очистить истёкшие токены: ", err)
	}

	return ct.RowsAffected(), nil
}
