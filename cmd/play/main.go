package main

import (
	"os"

	"log/slog"

	"github.com/hnasser-dev/wordle-solver/internal/game"
	"github.com/hnasser-dev/wordle-solver/internal/words"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	wordList, err := words.GetWordList()
	if err != nil {
		slog.Error("unable to read word list", "err", err)
		os.Exit(1)
	}

	freqMap, err := words.GetWordFrequencyMap()
	if err != nil {
		slog.Error("unable to read frequency map", "err", err)
		os.Exit(1)
	}

	answer := "hello"
	slog.Debug("answer", "answer", answer)

	guesses, gameWon := game.PlayGame(answer, wordList, freqMap)

	slog.Info(
		"game complete",
		slog.String("answer", answer),
		slog.Int("numGuesses", len(guesses)),
		slog.Any("guesses", guesses),
		slog.Bool("gameWon", gameWon),
	)
}
