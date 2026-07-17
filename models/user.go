package models

import "time"

type User struct {
	ID              string    `json:"id"`
	Email           string    `json:"email"`
	PasswordHash    string    `json:"-"`
	Username        string    `json:"username"`
	Rating          float64   `json:"rating"`
	RatingDeviation float64   `json:"rating_deviation"`
	CreatedAt       time.Time `json:"created_at"`
}
