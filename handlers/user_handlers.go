package handlers

import (
	"chess-puzzles-api/auth"
	"chess-puzzles-api/store"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func HandleGetProfile(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, ok := ctx.Value(auth.UserIDKey).(string)
		if !ok {
			slog.Error("invalid user id")
			http.Error(w, "internal server error", http.StatusUnauthorized)
			return
		}
		user, err := store.GetUserByID(ctx, pool, userID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				slog.Error("no user found with the given id:", "id", userID, "error", err)
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}
			slog.Error("internal server error:", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return

		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			slog.Error("internal server error", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

	}
}
