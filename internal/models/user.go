package models

import "time"

type User struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`       // NOT NULL
	Surname   string    `json:"surname"`    // NOT NULL
	Gender    string    `json:"gender"`     // NOT NULL
	BirthDate time.Time `json:"birth_date"` // NOT NULL

	HeightCm *int `json:"height_cm"` // nullable
	WeightKg *int `json:"weight_kg"` // nullable

	SportActivityLevelID     *int `json:"sport_activity_level_id"`     // FK nullable
	SportTargetID            *int `json:"sport_target_id"`             // FK nullable
	LocationPreferenceTypeID *int `json:"location_preference_type_id"` // FK nullable
	TownID                   *int `json:"town_id"`                     // FK nullable
	RoleID                   *int `json:"role_id"`                     // FK nullable

	PhoneNumber     string `json:"phone_number"`      // NOT NULL, UNIQUE
	IsPhoneVerified bool   `json:"is_phone_verified"` // NOT NULL, DEFAULT false
	Email           string `json:"email"`             // NOT NULL, UNIQUE
	IsEmailVerified bool   `json:"is_email_verified"` // NOT NULL, DEFAULT false
	Password        string `json:"-"`                 // hashed password

	IsHaveInjury      bool    `json:"is_have_injury"`     // DEFAULT false
	InjuryDescription *string `json:"injury_description"` // nullable
	Photo             *string `json:"photo"`              // nullable

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"` // nullable for soft delete
}
