package store

import (
	"chess-puzzles-api/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateSolveAttempt(ctx context.Context, pool *pgxpool.Pool, attempt models.SolveAttempt) error {
	query := `INSERT INTO solve_attempts (user_id,puzzle_id,success,user_rating_before,user_rating_after,puzzle_rating_before,puzzle_rating_after,created_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := pool.Exec(ctx, query,
		attempt.UserID,
		attempt.PuzzleID,
		attempt.Success,
		attempt.UserRatingBefore,
		attempt.UserRatingAfter,
		attempt.PuzzleRatingBefore,
		attempt.PuzzleRatingAfter,
		attempt.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("error inserting user: %w", err)
	}
	return nil

}
