package router

import (
	"chess-puzzles-api/auth"
	"chess-puzzles-api/config"
	"chess-puzzles-api/handlers"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRouter(pool *pgxpool.Pool, cfg *config.EnvConfig) *http.ServeMux {
	mux := http.NewServeMux()
	//Public Routes
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprint(w, "This is the Chess Puzzle API!")
		if err != nil {
			slog.Info("Error printing the http api print statement!")
		}
	})
	mux.HandleFunc("POST /auth/register", handlers.HandleRegister(pool, cfg))
	mux.HandleFunc("POST /auth/login", handlers.HandleLogin(pool, cfg))
	//Protected Routes
	mux.Handle("GET /puzzles/{id}", auth.AuthMiddleware(cfg.JWTSecret)(handlers.HandleGetPuzzle(pool)))
	mux.Handle("POST /puzzles/{id}/solve", auth.AuthMiddleware(cfg.JWTSecret)(handlers.HandleSolvePuzzle(pool)))
	mux.Handle("GET /users/me", auth.AuthMiddleware(cfg.JWTSecret)(handlers.HandleGetProfile(pool)))

	return mux
}
