package main

import (
	"chess-puzzles-api/auth"
	"chess-puzzles-api/config"
	"chess-puzzles-api/handlers"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRouter(pool *pgxpool.Pool, cfg *config.EnvConfig) *http.ServeMux {
	mux := http.NewServeMux()
	//Public Routes
	mux.HandleFunc("POST /auth/register", handlers.HandleRegister(pool, cfg))
	mux.HandleFunc("POST /auth/login", handlers.HandleLogin(pool, cfg))
	//Protected Routes
	mux.Handle("GET /puzzles/{id}", auth.AuthMiddleware(cfg.JWTSecret)(handlers.HandleGetPuzzle(pool)))
	mux.Handle("POST /puzzles/{id}/solve", auth.AuthMiddleware(cfg.JWTSecret)(handlers.HandleSolvePuzzle(pool)))
	mux.Handle("GET /users/me", auth.AuthMiddleware(cfg.JWTSecret)(handlers.HandleGetProfile(pool)))

	return mux
}
