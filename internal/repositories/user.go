package repositories

import (
	"context"
	"sport-assistance/internal/models"
	"sport-assistance/internal/services/dto"
)

func (r *Repository) CreateUser(ctx context.Context, user models.User) (uint64, error) {
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
		role_id,
		phone_number,
		email,
		password,
		is_have_injury,
		injury_description,
		photo
	)
	VALUES (
		$1,$2,$3,$4,
		$5,$6,$7,$8,COALESCE($9, (SELECT id FROM roles WHERE name = 'guest')),
		$10,$11,$12,$13,$14,$15
	)
	RETURNING id
	`

	var id uint64
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
		user.RoleID,
		user.PhoneNumber,
		user.Email,
		user.Password,
		user.IsHaveInjury,
		user.InjuryDescription,
		user.Photo,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) GetUserByID(ctx context.Context, userID uint64) (dto.UserDto, error) {
	query := `
		SELECT
			id,
			name,
			surname,
			gender,
			birth_date,
			height_cm,
			weight_kg,
			sport_activity_level_id,
			town_id,
			role_id,
			phone_number,
			email,
			password,
			is_have_injury,
			injury_description,
			photo,
			created_at,
			updated_at,
			deleted_at
		FROM users
		WHERE id = $1
		`

	var user models.User
	err := r.postgres.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Surname,
		&user.Gender,
		&user.BirthDate,
		&user.HeightCm,
		&user.WeightKg,
		&user.SportActivityLevelID,
		&user.TownID,
		&user.RoleID,
		&user.PhoneNumber,
		&user.Email,
		&user.Password,
		&user.IsHaveInjury,
		&user.InjuryDescription,
		&user.Photo,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		return dto.UserDto{}, err
	}

	userDTO, err := dto.UserToDto(user)
	if err != nil {
		return dto.UserDto{}, err
	}

	return userDTO, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (dto.UserDto, error) {
	q := `SELECT id, email, password FROM users WHERE email = $1`
	var user models.User
	err := r.postgres.QueryRow(ctx, q, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return dto.UserDto{}, err
	}

	userDTO, err := dto.UserToDto(user)
	if err != nil {
		return dto.UserDto{}, err
	}

	return userDTO, nil
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
