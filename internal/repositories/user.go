package repositories

import (
	"context"
	"sport-assistance/internal/models"
	"sport-assistance/internal/services/dto"

	"github.com/jackc/pgx/v5"
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
		sport_target_id,
		location_preference_type_id,
		town_id,
		role_id,
		phone_number,
		is_phone_verified,
		email,
		is_email_verified,
		password,
		is_have_injury,
		injury_description,
		photo
	)
	VALUES (
		$1,$2,$3,$4,
		$5,$6,$7,$8,$9,$10,COALESCE($11, (SELECT id FROM roles WHERE name = 'guest')),
		$12,$13,$14,$15,$16,$17,$18,$19
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
		user.SportTargetID,
		user.LocationPreferenceTypeID,
		user.TownID,
		user.RoleID,
		user.PhoneNumber,
		user.IsPhoneVerified,
		user.Email,
		user.IsEmailVerified,
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
			sport_target_id,
			location_preference_type_id,
			town_id,
			role_id,
			phone_number,
			is_phone_verified,
			email,
			is_email_verified,
			password,
			is_have_injury,
			injury_description,
			photo,
			created_at,
			updated_at,
			deleted_at
		FROM users
		WHERE id = $1
		  AND deleted_at IS NULL
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
		&user.SportTargetID,
		&user.LocationPreferenceTypeID,
		&user.TownID,
		&user.RoleID,
		&user.PhoneNumber,
		&user.IsPhoneVerified,
		&user.Email,
		&user.IsEmailVerified,
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
	q := `SELECT id, email, password FROM users WHERE email = $1 AND deleted_at IS NULL`
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

func (r *Repository) GetUserByPhone(ctx context.Context, phone string) (dto.UserDto, error) {
	q := `SELECT id, email FROM users WHERE phone_number = $1 AND deleted_at IS NULL`
	var user models.User
	err := r.postgres.QueryRow(ctx, q, phone).Scan(&user.ID, &user.Email)
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

func (r *Repository) GetUsers(ctx context.Context) ([]dto.UserDto, error) {
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
			sport_target_id,
			location_preference_type_id,
			town_id,
			role_id,
			phone_number,
			is_phone_verified,
			email,
			is_email_verified,
			password,
			is_have_injury,
			injury_description,
			photo,
			created_at,
			updated_at,
			deleted_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY id
	`

	rows, err := r.postgres.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]dto.UserDto, 0)
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Surname,
			&user.Gender,
			&user.BirthDate,
			&user.HeightCm,
			&user.WeightKg,
			&user.SportActivityLevelID,
			&user.SportTargetID,
			&user.LocationPreferenceTypeID,
			&user.TownID,
			&user.RoleID,
			&user.PhoneNumber,
			&user.IsPhoneVerified,
			&user.Email,
			&user.IsEmailVerified,
			&user.Password,
			&user.IsHaveInjury,
			&user.InjuryDescription,
			&user.Photo,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		); err != nil {
			return nil, err
		}

		userDTO, err := dto.UserToDto(user)
		if err != nil {
			return nil, err
		}

		users = append(users, userDTO)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) UpdateUser(ctx context.Context, userID uint64, user models.User) error {
	query := `
		UPDATE users
		SET
			name = $2,
			surname = $3,
			gender = $4,
			birth_date = $5,
			height_cm = $6,
			weight_kg = $7,
			sport_activity_level_id = $8,
			sport_target_id = $9,
			location_preference_type_id = $10,
			town_id = $11,
			role_id = $12,
			phone_number = $13,
			is_phone_verified = $14,
			email = $15,
			is_email_verified = $16,
			password = $17,
			is_have_injury = $18,
			injury_description = $19,
			photo = $20,
			updated_at = now()
		WHERE id = $1
		  AND deleted_at IS NULL
	`

	ct, err := r.postgres.Exec(
		ctx,
		query,
		userID,
		user.Name,
		user.Surname,
		user.Gender,
		user.BirthDate,
		user.HeightCm,
		user.WeightKg,
		user.SportActivityLevelID,
		user.SportTargetID,
		user.LocationPreferenceTypeID,
		user.TownID,
		user.RoleID,
		user.PhoneNumber,
		user.IsPhoneVerified,
		user.Email,
		user.IsEmailVerified,
		user.Password,
		user.IsHaveInjury,
		user.InjuryDescription,
		user.Photo,
	)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, userID uint64) error {
	query := `
		UPDATE users
		SET deleted_at = now(), updated_at = now()
		WHERE id = $1
		  AND deleted_at IS NULL
	`

	ct, err := r.postgres.Exec(ctx, query, userID)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
