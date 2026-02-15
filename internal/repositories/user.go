package repositories

import (
	"context"
	"sport-assistance/internal/models"
)

func (r *Repository) CreateUser(ctx context.Context, user models.User) (int64, error) {
	query := `
	INSERT INTO users (
		name,
		surname,
		gender,
		birth_date,
		height_cm,
		weight_kg,
		sport_activity_level_id,
		town_id,
		phone_number,
		email,
		password,
		is_have_injury,
		injury_description,
		photo
	)
	VALUES (
		$1,$2,$3,$4,
		$5,$6,$7,$8,
		$9,$10,$11,$12,$13,$14
	)
	RETURNING id
	`

	var id int64
	err := r.postgres.QueryRow(
		ctx,
		query,
		user.Name,
		user.Surname,
		user.Gender,
		user.BirthDate,
		user.HeightCm,
		user.WeightKg,
		user.SportActivityLevelID,
		user.TownID,
		user.PhoneNumber,
		user.Email,
		user.Password,
		user.IsHaveInjury,
		user.InjuryDescription,
		user.Photo,
	).Scan(&id)

	if err != nil {
		return -1, err
	}

	return id, nil
}

func (r *Repository) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `
	SELECT EXISTS (
		SELECT 1
		FROM users
		WHERE email = $1
		  AND deleted_at IS NULL
	)
	`

	var exists bool
	if err := r.postgres.QueryRow(ctx, query, email).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
