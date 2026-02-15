package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`       // NOT NULL
	Surname   string    `json:"surname"`    // NOT NULL
	Gender    string    `json:"gender"`     // NOT NULL
	BirthDate time.Time `json:"birth_date"` // NOT NULL

	HeightCm *int `json:"height_cm"` // nullable
	WeightKg *int `json:"weight_kg"` // nullable

	SportActivityLevelID *int `json:"sport_activity_level_id"` // FK nullable
	TownID               *int `json:"town_id"`                 // FK nullable

	PhoneNumber string `json:"phone_number"` // NOT NULL, UNIQUE
	Email       string `json:"email"`        // NOT NULL, UNIQUE
	Password    string `json:"-"`            // хешированный пароль

	IsHaveInjury      bool    `json:"is_have_injury"`     // DEFAULT false
	InjuryDescription *string `json:"injury_description"` // nullable
	Photo             *string `json:"photo"`              // nullable

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"` // nullable для soft delete
}
