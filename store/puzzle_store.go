package store

import (
	"chess-puzzles-api/models"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrDuplicatePuzzle = errors.New("puzzle already exists")
var ErrInvalidReference = errors.New("referenced item not found")

func GetPuzzlebyID(ctx context.Context, pool *pgxpool.Pool, id string) (models.Puzzle, error) {
	query := `SELECT id,fen,moves,rating,rating_deviation,popularity,nb_plays,themes,created_at FROM puzzles WHERE id=$1`
	var p models.Puzzle
	err := pool.QueryRow(ctx, query, id).Scan(
		&p.ID,
		&p.FEN,
		&p.Moves,
		&p.Rating,
		&p.RatingDeviation,
		&p.Popularity,
		&p.NbPlays,
		&p.Themes,
		&p.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Puzzle{}, fmt.Errorf("puzzle %s not found: %w", id, err)
		}
		return models.Puzzle{}, err

	}
	return p, nil
}

func CreatePuzzle(ctx context.Context, pool *pgxpool.Pool, p models.Puzzle) error {
	query := `INSERT INTO puzzles (id, fen,moves,rating,rating_deviation,popularity,nb_plays,themes,created_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`

	_, err := pool.Exec(ctx, query, p.ID, p.FEN, p.Moves, p.Rating, p.RatingDeviation, p.Popularity, p.NbPlays, p.Themes, p.CreatedAt)
	if err != nil {
		//to handle specific postgres engine violations
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return fmt.Errorf("bad request: record already exists (%w)", ErrDuplicatePuzzle)

			case "23503":
				return fmt.Errorf("bad request: referenced parent item not found: %w", ErrInvalidReference)
			}

		}
		return fmt.Errorf("unexpected database error: %w", err)
	}
	return nil
}
func UpdatePuzzleRating(ctx context.Context, pool *pgxpool.Pool, puzzleID string, newRating, newRD float64) error {
	query := `UPDATE puzzles SET rating=$1,rating_deviation=$2 WHERE id=$3;`
	commandTag, err := pool.Exec(ctx, query, newRating, newRD, puzzleID)
	if err != nil {
		return fmt.Errorf("failed to update puzzle rating: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no puzzle found with ID %s", puzzleID)
	}
	return nil

}
func GetRandomPuzzle(ctx context.Context, pool *pgxpool.Pool) (models.Puzzle, error) {
	query := `SELECT id,fen,moves,rating,rating_deviation,popularity,nb_plays,themes,created_at 
			FROM puzzles ORDER BY RANDOM() LIMIT 1`
	var p models.Puzzle
	err := pool.QueryRow(ctx, query).Scan(
		&p.ID,
		&p.FEN,
		&p.Moves,
		&p.Rating,
		&p.RatingDeviation,
		&p.Popularity,
		&p.NbPlays,
		&p.Themes,
		&p.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Puzzle{}, fmt.Errorf("puzzle not found: %w", err)
		}
		return models.Puzzle{}, err

	}
	return p, nil

}
