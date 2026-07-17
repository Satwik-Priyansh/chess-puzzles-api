package main

import (
	"chess-puzzles-api/config"
	"chess-puzzles-api/db"
	"fmt"
	"log/slog"
	"net/http"
	"os"
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
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprint(w, "This is the Chess Puzzle API!")
		if err != nil {
			slog.Info("Error printing the http api print statement!")
		}
	})
	slog.Info("Server listening on port 3000")
	err_server := http.ListenAndServe(":3000", mux)
	if err_server != nil {
		slog.Error("Failed to start http server", "error", err_server)
		os.Exit(1)
	}

}
