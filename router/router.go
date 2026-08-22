package router

import (
	"chess-puzzles-api/auth"
	"chess-puzzles-api/config"
	"chess-puzzles-api/handlers"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/time/rate"
)

func SetupRouter(pool *pgxpool.Pool, cfg *config.EnvConfig) http.Handler {
	mux := http.NewServeMux()
	//Public Routes
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprint(w, "This is the Chess Puzzle API!")
		if err != nil {
			slog.Info("Error printing the http api print statement!")
		}
	})
	authLimiter := auth.NewIPRateLimiter(rate.Every(time.Minute), 10)
	mux.Handle("POST /auth/register", auth.RateLimitMiddleware(authLimiter)(handlers.HandleRegister(pool, cfg)))
	mux.Handle("POST /auth/login", auth.RateLimitMiddleware(authLimiter)(handlers.HandleLogin(pool, cfg)))
	mux.HandleFunc("GET /leaderboard", handlers.HandleGetLeaderboard(pool))
	//Protected Routes
	mux.Handle("GET /puzzles/random", auth.AuthMiddleware(cfg.JWTSecret)(handlers.HandleGetRandomPuzzle(pool)))
	mux.Handle("GET /puzzles/{id}", auth.AuthMiddleware(cfg.JWTSecret)(handlers.HandleGetPuzzle(pool)))
	mux.Handle("POST /puzzles/{id}/solve", auth.AuthMiddleware(cfg.JWTSecret)(handlers.HandleSolvePuzzle(pool)))
	mux.Handle("GET /users/me", auth.AuthMiddleware(cfg.JWTSecret)(handlers.HandleGetProfile(pool)))

	return auth.CORSMiddleware(auth.SecurityHeaders(auth.RequestLogger(auth.MaxBodySize(mux))))
}
