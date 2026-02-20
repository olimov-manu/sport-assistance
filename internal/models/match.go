package models

import "time"

type Match struct {
	ID          uint64    `json:"id"`
	MatchTypeID int       `json:"match_type_id"`
	CreatedAt   time.Time `json:"created_at"`
}
