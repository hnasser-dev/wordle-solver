package main

import (
	"os"

	"log/slog"

	"github.com/hnasser-dev/wordle-solver/internal/game"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	answer := "hello"
	slog.Debug("answer", "answer", answer)

	// initialGuesses := []string{"hello"}
	// game, err := game.NewGame(answer, initialGuesses...)

	game, err := game.NewGame(answer)
	if err != nil {
		slog.Error("unable to create new game", "err", err)
		os.Exit(1)
	}

	gameWon, guesses := game.PlayGameUntilEnd(true)

	slog.Info(
		"game complete",
		slog.String("answer", answer),
		slog.Int("numGuesses", len(guesses)),
		slog.Any("guesses", guesses),
		slog.Bool("gameWon", gameWon),
	)
}
