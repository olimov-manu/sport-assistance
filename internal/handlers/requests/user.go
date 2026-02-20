package requests

type CreateUserRequest struct {
	Name      string `json:"name" validate:"required,max=100"`
	Surname   string `json:"surname" validate:"required,max=100"`
	Gender    string `json:"gender" validate:"required"`
	BirthDate string `json:"birth_date" validate:"required,datetime=02-01-2006"` // Format: DD-MM-YYYY

	HeightCm *int `json:"height_cm,omitempty"`
	WeightKg *int `json:"weight_kg,omitempty"`

	SportActivityLevelID *int `json:"sport_activity_level_id,omitempty"`
	TownID               *int `json:"town_id,omitempty"`
	RoleID               *int `json:"role_id,omitempty"`

	PhoneNumber string `json:"phone_number" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=6"`

	IsHaveInjury      bool    `json:"is_have_injury"`
	InjuryDescription *string `json:"injury_description,omitempty"`
	Photo             *string `json:"photo,omitempty"`
}
