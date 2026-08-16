package store

import (
	"chess-puzzles-api/models"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(ctx context.Context, pool *pgxpool.Pool, user models.User) (string, error) {
	query := `INSERT INTO users (email,password_hash,username,rating,rating_deviation,created_at)
	VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`
	var id string
	err := pool.QueryRow(ctx, query, user.Email, user.PasswordHash, user.Username, user.Rating, user.RatingDeviation, user.CreatedAt).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("error inserting user: %w", err)
	}
	return id, nil
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
func GetUserByID(ctx context.Context, pool *pgxpool.Pool, id string) (models.User, error) {
	query := `SELECT *  FROM users WHERE id=$1;`
	var fetchedUser models.User
	err := pool.QueryRow(ctx, query, id).Scan(
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
			return fetchedUser, fmt.Errorf("user lookup complete: no user found with id: %s ", id)
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
func GetTopUsers(ctx context.Context, pool *pgxpool.Pool, limit int) ([]models.User, error) {
	query := `SELECT id,email,username,rating,rating_deviation,created_at 
FROM users ORDER BY rating DESC LIMIT $1`
	rows, err := pool.Query(ctx, query, limit)
	if err != nil {
		slog.Error("error fetching multiple rows", "error", err)
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Email, &u.Username, &u.Rating, &u.RatingDeviation, &u.CreatedAt)
		if err != nil {
			slog.Error("error while scanning users.", "error", err)
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}
	return users, nil

}
