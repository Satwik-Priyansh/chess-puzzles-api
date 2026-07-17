package store

import (
	"chess-puzzles-api/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)


func CreateUser (ctx context.Context, pool *pgxpool.Pool, user models.User) error{
	query:= `INSERT INTO users (user)`
}