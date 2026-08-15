package store

import (
	"chess-puzzles-api/models"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(ctx context.Context, pool *pgxpool.Pool, user models.User) error {
	query := `INSERT INTO users (email,password_hash,username,rating,rating_deviation,created_at)
	values($1,$2,$3,$4,$5,$6);`
	_, err := pool.Exec(ctx, query, user.Email, user.PasswordHash, user.Username, user.Rating, user.RatingDeviation, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("error inserting row %v.", err)
	}
	return nil
}
func GetUserByEmail(ctx context.Context, pool *pgxpool.Pool, email string) (models.User, error) {
	query := `SELECT *  FROM users WHERE email=$1;`
	var fetchedUser models.User
	err := pool.QueryRow(ctx, query, email).Scan(
		&fetchedUser.ID,
		&fetchedUser.Email,
		&fetchedUser.PasswordHash,
		&fetchedUser.Username,
		&fetchedUser.Rating,
		&fetchedUser.RatingDeviation,
		&fetchedUser.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fetchedUser, fmt.Errorf("user loookup complete: no user found with email %s ", email)
		}
		return fetchedUser, fmt.Errorf("database error:%v", err)
	}
	return fetchedUser, nil
}
func UpdateUserRating(ctx context.Context, pool *pgxpool.Pool, userID string, newRating, newRD float64) error {
	query := `UPDATE users SET rating=$1,rating_deviation=$2 WHERE id=$3;`
	commandTag, err := pool.Exec(ctx, query, newRating, newRD, userID)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no user found with ID %s", userID)
	}
	return nil

}
