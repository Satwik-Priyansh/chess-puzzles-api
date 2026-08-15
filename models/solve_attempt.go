package models

import "time"

type SolveAttempt struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	PuzzleID           string    `json:"puzzle_id"`
	Success            bool      `json:"success"`
	UserRatingBefore   float64   `json:"user_rating_before"`
	UserRatingAfter    float64   `json:"user_rating_after"`
	PuzzleRatingBefore float64   `json:"puzzle_rating_before"`
	PuzzleRatingAfter  float64   `json:"puzzle_rating_after"`
	CreatedAt          time.Time `json:"created_at"`
}
