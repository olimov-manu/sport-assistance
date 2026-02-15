package models

import "time"

type RefreshToken struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	Token     string    `json:"token"`
	Expire    int       `json:"expire"`
	CreatedAt time.Time `json:"created_at"`
	RevokedAT time.Time `json:"revoked_at"`
}
