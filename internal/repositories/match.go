package repositories

import (
	"context"
	"sport-assistance/internal/models"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateMatch(ctx context.Context, matchTypeID int, participantIDs []uint64) (uint64, error) {
	tx, err := r.postgres.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	const insertMatchQuery = `
		INSERT INTO matches (match_type_id)
		VALUES ($1)
		RETURNING id
	`

	var matchID uint64
	if err = tx.QueryRow(ctx, insertMatchQuery, matchTypeID).Scan(&matchID); err != nil {
		return 0, err
	}

	if len(participantIDs) > 0 {
		const insertParticipantQuery = `
			INSERT INTO user_matches (user_id, match_id)
			VALUES ($1, $2)
			ON CONFLICT (user_id, match_id) DO NOTHING
		`

		for _, participantID := range participantIDs {
			if _, err = tx.Exec(ctx, insertParticipantQuery, participantID, matchID); err != nil {
				return 0, err
			}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}

	return matchID, nil
}

func (r *Repository) GetMatchByID(ctx context.Context, matchID uint64) (models.Match, error) {
	const query = `
		SELECT id, match_type_id, created_at
		FROM matches
		WHERE id = $1
	`

	var match models.Match
	if err := r.postgres.QueryRow(ctx, query, matchID).Scan(&match.ID, &match.MatchTypeID, &match.CreatedAt); err != nil {
		return models.Match{}, err
	}

	return match, nil
}

func (r *Repository) GetMatchesByUserID(ctx context.Context, userID uint64) ([]models.Match, error) {
	const query = `
		SELECT m.id, m.match_type_id, m.created_at
		FROM matches m
		JOIN user_matches um ON um.match_id = m.id
		WHERE um.user_id = $1
		ORDER BY m.created_at DESC, m.id DESC
	`

	rows, err := r.postgres.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make([]models.Match, 0)
	for rows.Next() {
		var match models.Match
		if err := rows.Scan(&match.ID, &match.MatchTypeID, &match.CreatedAt); err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

func (r *Repository) GetMatchParticipants(ctx context.Context, matchID uint64) ([]uint64, error) {
	const query = `
		SELECT user_id
		FROM user_matches
		WHERE match_id = $1
		ORDER BY user_id
	`

	rows, err := r.postgres.Query(ctx, query, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	participants := make([]uint64, 0)
	for rows.Next() {
		var userID uint64
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		participants = append(participants, userID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return participants, nil
}

func (r *Repository) AddUserToMatch(ctx context.Context, matchID, userID uint64) error {
	const query = `
		INSERT INTO user_matches (user_id, match_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, match_id) DO NOTHING
	`

	_, err := r.postgres.Exec(ctx, query, userID, matchID)
	return err
}

func (r *Repository) RemoveUserFromMatch(ctx context.Context, matchID, userID uint64) error {
	const query = `
		DELETE FROM user_matches
		WHERE user_id = $1
		  AND match_id = $2
	`

	ct, err := r.postgres.Exec(ctx, query, userID, matchID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *Repository) IsUserInMatch(ctx context.Context, matchID, userID uint64) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM user_matches
			WHERE user_id = $1
			  AND match_id = $2
		)
	`

	var exists bool
	if err := r.postgres.QueryRow(ctx, query, userID, matchID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (r *Repository) UpdateMatchType(ctx context.Context, matchID uint64, matchTypeID int) error {
	const query = `
		UPDATE matches
		SET match_type_id = $2
		WHERE id = $1
	`

	ct, err := r.postgres.Exec(ctx, query, matchID, matchTypeID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *Repository) DeleteMatch(ctx context.Context, matchID uint64) error {
	const query = `DELETE FROM matches WHERE id = $1`

	ct, err := r.postgres.Exec(ctx, query, matchID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
