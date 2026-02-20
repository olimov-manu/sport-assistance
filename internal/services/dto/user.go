package dto

import (
	"sport-assistance/internal/models"
	"time"
)

type UserDto struct {
	ID                       uint64    `json:"id"`
	Name                     string    `json:"name"`
	Surname                  string    `json:"surname"`
	Gender                   string    `json:"gender"`
	BirthDate                time.Time `json:"birth_date"`
	HeightCm                 *int      `json:"height_cm"`
	WeightKg                 *int      `json:"weight_kg"`
	SportActivityLevelID     *int      `json:"sport_activity_level_id"`
	SportTargetID            *int      `json:"sport_target_id"`
	LocationPreferenceTypeID *int      `json:"location_preference_type_id"`
	TownID                   *int      `json:"town_id"`
	RoleID                   *int      `json:"role_id"`
	PhoneNumber              string    `json:"phone_number"`
	IsPhoneVerified          bool      `json:"is_phone_verified"`
	Email                    string    `json:"email"`
	IsEmailVerified          bool      `json:"is_email_verified"`
	Password                 string    `json:"password"`
	IsHaveInjury             bool      `json:"is_have_injury"`
	InjuryDescription        *string   `json:"injury_description"`
	Photo                    *string   `json:"photo"`
}

func UserToDto(userModel models.User) (UserDto, error) {
	return UserDto{
		ID:                       userModel.ID,
		Name:                     userModel.Name,
		Surname:                  userModel.Surname,
		Gender:                   userModel.Gender,
		BirthDate:                userModel.BirthDate,
		HeightCm:                 userModel.HeightCm,
		WeightKg:                 userModel.WeightKg,
		SportActivityLevelID:     userModel.SportActivityLevelID,
		SportTargetID:            userModel.SportTargetID,
		LocationPreferenceTypeID: userModel.LocationPreferenceTypeID,
		TownID:                   userModel.TownID,
		RoleID:                   userModel.RoleID,
		PhoneNumber:              userModel.PhoneNumber,
		IsPhoneVerified:          userModel.IsPhoneVerified,
		Email:                    userModel.Email,
		IsEmailVerified:          userModel.IsEmailVerified,
		Password:                 userModel.Password,
		IsHaveInjury:             userModel.IsHaveInjury,
		InjuryDescription:        userModel.InjuryDescription,
		Photo:                    userModel.Photo,
	}, nil
}
