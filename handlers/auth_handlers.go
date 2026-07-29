package handlers

import (
	"chess-puzzles-api/auth"
	"chess-puzzles-api/config"
	"chess-puzzles-api/models"
	"chess-puzzles-api/store"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func HandleRegister(pool *pgxpool.Pool, cfg *config.EnvConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Username string `json:"username"` //never declare this on package level, it will cause race condition if two user send post request simultaneously
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid JSON payload:"+err.Error(), http.StatusBadRequest)
			return
		}
		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			slog.Error("error hashing the password:", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		user := &models.User{Rating: 1000.0, RatingDeviation: 350.0, PasswordHash: hashedPassword, Email: req.Email, Username: req.Username}
		err = store.CreateUser(r.Context(), pool, *user)
		if err != nil {
			slog.Error("error creating the user:", "error", err)
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			slog.Error("error while encoding:", "error", err)
			return
		}

	}
}
