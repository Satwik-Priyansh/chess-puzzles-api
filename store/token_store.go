package store

import (
	"chess-puzzles-api/auth"
	"chess-puzzles-api/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateRefreshToken(ctx context.Context, pool *pgxpool.Pool, userID string, expiresAt time.Time) (string, error) {
	refreshToken, err := auth.GenerateSecureToken()
	if err != nil {
		return "", err
	}
	query := `INSERT INTO refresh_tokens (user_id,token,expires_at) VALUES ($1,$2,$3)`
	_, err = pool.Exec(ctx, query, userID, refreshToken, expiresAt)
	if err != nil {
		return "", fmt.Errorf("error inserting refresh token: %w", err)
	}
	return refreshToken, nil

}
func GetRefreshToken(ctx context.Context, pool *pgxpool.Pool, token string) (models.RefreshToken, error) {
	var fetchedRefreshToken models.RefreshToken
	query := `SELECT * FROM refresh_tokens WHERE token = $1`
	err := pool.QueryRow(ctx, query, token).Scan(
		&fetchedRefreshToken.ID,
		&fetchedRefreshToken.UserID,
		&fetchedRefreshToken.Token,
		&fetchedRefreshToken.ExpiresAt,
		&fetchedRefreshToken.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fetchedRefreshToken, fmt.Errorf("token lookup complete: no token found")
		}
		return fetchedRefreshToken, fmt.Errorf("database error:%v", err)
	}
	return fetchedRefreshToken, nil
}
func DeleteRefreshToken(ctx context.Context, pool *pgxpool.Pool, token string) error {
	query := `DELETE FROM refresh_tokens WHERE token= $1`
	commandTag, err := pool.Exec(ctx, query, token)
	if err != nil {
		return fmt.Errorf("error while deleting token:%v", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("token: %v not found", token)
	}

	return nil

}
func DeleteAllUserTokens(ctx context.Context, pool *pgxpool.Pool, userID string) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	commandTag, err := pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("error while deleting all tokens with assocaited user id:%v", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no token found with userID:%s", userID)
	}
	return nil
} // for logout-all-devices
