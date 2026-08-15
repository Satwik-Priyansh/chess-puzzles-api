package handlers

import (
	"chess-puzzles-api/store"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func HandleGetPuzzle(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		puzzle, err := store.GetPuzzlebyID(r.Context(), pool, id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				slog.Error("error finding puzzle", "error", err)
				http.Error(w, "puzzle not found"+err.Error(), http.StatusNotFound)
				return
			}
			slog.Error("internal server error", "error", err)
			http.Error(w, "puzzle not found"+err.Error(), http.StatusInternalServerError)
			return

		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(puzzle)
		if err != nil {
			slog.Error("error encoding puzzle", "error", err)
			http.Error(w, "internal server error"+err.Error(), http.StatusInternalServerError)
		}

	}
}
