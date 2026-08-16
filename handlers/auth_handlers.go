package handlers

import (
	"chess-puzzles-api/auth"
	"chess-puzzles-api/config"
	"chess-puzzles-api/models"
	"chess-puzzles-api/store"
	"chess-puzzles-api/validation"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenResponse struct {
	Token string `json:"token"`
}

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
		if err := validation.ValidateEmail(req.Email); err != nil {
			http.Error(w, "invalid email format", http.StatusBadRequest)
			return
		}
		if err := validation.ValidateUsername(req.Username); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := validation.ValidatePassword(req.Password); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			slog.Error("error hashing the password:", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		user := &models.User{Rating: 1000.0, RatingDeviation: 350.0, PasswordHash: hashedPassword, Email: req.Email, Username: req.Username}
		userID, err := store.CreateUser(r.Context(), pool, *user)
		if err != nil {
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}
		user.ID = userID
		token, err := auth.GenerateToken(userID, cfg.JWTSecret)
		if err != nil {
			http.Error(w, "error generating token", http.StatusInternalServerError)
			return
		}
		response := struct {
			Token string      `json:"token"`
			User  models.User `json:"user"`
		}{
			Token: token,
			User:  *user,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			slog.Error("error while encoding:", "error", err)
			return
		}

	}
}

func HandleLogin(pool *pgxpool.Pool, cfg *config.EnvConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid JSON payload:"+err.Error(), http.StatusBadRequest)
			return
		}
		user, err := store.GetUserByEmail(r.Context(), pool, req.Email)
		if err != nil {
			slog.Error("user not found", "error", err)
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}
		err = auth.CheckPassword(req.Password, user.PasswordHash)
		if err != nil {
			slog.Error("wrong password:", "error", err)
			http.Error(w, "incorrect username/password", http.StatusUnauthorized)
			return
		}
		token, err := auth.GenerateToken(user.ID, cfg.JWTSecret)
		if err != nil {
			slog.Error("error generating token:", "error", err)
			http.Error(w, "authentication error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		response := TokenResponse{
			Token: token,
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			slog.Error("error encoding response", "error", err)
		}

	}

}
