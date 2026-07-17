package models

import "time"

type Puzzle struct {
	ID              string    `json:"id"` // text id, e.g. "0000D"
	FEN             string    `json:"fen"`
	Moves           []string  `json:"moves"` // split from space-separated UCI moves
	Rating          float64   `json:"rating"`
	RatingDeviation float64   `json:"rating_deviation"`
	Popularity      int       `json:"popularity"`
	NbPlays         int       `json:"nb_plays"`
	Themes          []string  `json:"themes"` // split from space-separated themes
	CreatedAt       time.Time `json:"created_at"`
}
