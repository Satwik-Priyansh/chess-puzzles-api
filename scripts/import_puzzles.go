package main

import (
	"chess-puzzles-api/config"
	"chess-puzzles-api/db"
	"chess-puzzles-api/models"
	"chess-puzzles-api/store"
	"context"
	"encoding/csv"
	"errors"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	cfg := config.LoadConfig()
	conn, err := db.ConnectDB(cfg)
	if err != nil {
		slog.Error("Database connection falied", "error", err)
		os.Exit(1)
	} else {
		slog.Info("Database connected successfully.")
	}
	defer conn.Close()
	file, err := os.Open("/Users/satwikpriyansh/Projects/Golang_Projects/chess-puzzle-crud-api/scripts/lichess_db_puzzle.csv")
	if err != nil {
		slog.Error("failed to open file", "error", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	_, err = reader.Read()
	if err != nil {
		slog.Error("Error while reading file ", "error", err)
	}
	const maxRows = 1000
	rowCount := 0
	for rowCount < maxRows {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		rowCount++
		if err != nil {
			slog.Error("Error reading", "row", rowCount+1, "error", err)
		}
		if len(record) < 8 {
			slog.Error("skipping malformed row", "row", record)
			continue
		}
		rating, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			slog.Error("Invalid rating field", "error", err)
			continue
		}
		rating_deviation, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			slog.Error("Invalid rating deviation field", "error", err)
			continue
		}
		popularity, err := strconv.Atoi(record[5])
		if err != nil {
			slog.Error("invalid field popularity", "error", err)
			continue
		}
		moves := strings.Split(record[2], " ")
		nb_plays, err := strconv.Atoi(record[6])
		if err != nil {
			slog.Error("Invalid field NbPlays", "error", err)
			continue
		}
		themes := strings.Split(record[7], " ")

		puzzle := models.Puzzle{
			ID:              record[0],
			FEN:             record[1],
			Moves:           moves,
			Rating:          rating,
			RatingDeviation: rating_deviation,
			Popularity:      popularity,
			NbPlays:         nb_plays,
			Themes:          themes,
			CreatedAt:       time.Now(),
		}
		err = store.CreatePuzzle(context.Background(), conn, puzzle)
		if err == nil {
			slog.Info("puzzle added successfully")
		} else if errors.Is(err, store.ErrDuplicatePuzzle) {
			slog.Info("found duplaicate puzzle skipping to next one")
			continue

		} else if errors.Is(err, store.ErrInvalidReference) {
			slog.Info("referenced item not found")
			continue
		} else {
			slog.Error("database error", "error", err)
			os.Exit(1)
		}
	}
	slog.Info("import completed with", "rows", rowCount)

}
