package handlers

import (
	"chess-puzzles-api/auth"
	"chess-puzzles-api/models"
	"chess-puzzles-api/rating"
	"chess-puzzles-api/store"
	"chess-puzzles-api/validation"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

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
				http.Error(w, "puzzle not found", http.StatusNotFound)
				return
			}
			slog.Error("internal server error", "error", err)
			http.Error(w, "puzzle not found", http.StatusInternalServerError)
			return

		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(puzzle)
		if err != nil {
			slog.Error("error encoding puzzle", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}

	}
}
func HandleSolvePuzzle(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, ok := ctx.Value(auth.UserIDKey).(string)
		if !ok {
			slog.Error("invalid user id")
			http.Error(w, "internal server error", http.StatusUnauthorized)
			return
		}
		puzzleID := r.PathValue("id")
		puzzle, err := store.GetPuzzlebyID(ctx, pool, puzzleID)
		if err != nil {
			slog.Error("puzzle not found", "error", err)
			http.Error(w, "puzzle not found", http.StatusNotFound)
			return

		}
		user, err := store.GetUserByID(ctx, pool, userID)
		if err != nil {
			slog.Error("user not found", "error", err)
			http.Error(w, "user not found", http.StatusNotFound)
			return

		}
		var req struct {
			Moves []string `json:"moves"`
		}
		err = json.NewDecoder(r.Body).Decode(&req)
		if err := validation.ValidateUCIMoves(req.Moves); err != nil {
			http.Error(w, "invalid moves", http.StatusBadRequest)
			return
		}
		if err != nil {
			slog.Error("error decoding json", "error", err)
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		solved := true
		puzzleIndex := 1
		for _, value := range req.Moves {
			if puzzleIndex >= len(puzzle.Moves) || value != puzzle.Moves[puzzleIndex] {
				solved = false
				break
			}
			puzzleIndex += 1
		}
		newUserRating, newUserRatingDev, newPuzzleRating, newPuzzleRatingDev := rating.CalculateNewRatings(user.Rating, user.RatingDeviation, puzzle.Rating, puzzle.RatingDeviation, solved)
		err = store.UpdateUserRating(ctx, pool, userID, newUserRating, newUserRatingDev)
		if err != nil {
			slog.Error("error updating user rating", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		err = store.UpdatePuzzleRating(ctx, pool, puzzleID, newPuzzleRating, newPuzzleRatingDev)
		if err != nil {
			slog.Error("error updating puzzle rating", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		attempt := models.SolveAttempt{
			UserID:             userID,
			PuzzleID:           puzzleID,
			Success:            solved,
			UserRatingBefore:   user.Rating,
			UserRatingAfter:    newUserRating,
			PuzzleRatingBefore: puzzle.Rating,
			PuzzleRatingAfter:  newPuzzleRating,
			CreatedAt:          time.Now(),
		}
		err = store.CreateSolveAttempt(ctx, pool, attempt)
		if err != nil {
			slog.Error("error creating solve attempt", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		type responseStruct struct {
			Success      bool    `json:"success"`
			NewRating    float64 `json:"new_rating"`
			RatingChange float64 `json:"rating_change"`
		}
		response := responseStruct{
			Success:      solved,
			NewRating:    newUserRating,
			RatingChange: newUserRating - user.Rating,
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			slog.Error("error encoding response", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}

	}
}
func HandleGetRandomPuzzle(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		puzzle, err := store.GetRandomPuzzle(r.Context(), pool)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				slog.Error("error finding puzzle", "error", err)
				http.Error(w, "puzzle not found", http.StatusNotFound)
				return
			}
			slog.Error("internal server error", "error", err)
			http.Error(w, "puzzle not found", http.StatusInternalServerError)
			return

		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(puzzle)
		if err != nil {
			slog.Error("error encoding puzzle", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}
