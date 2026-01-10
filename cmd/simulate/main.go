package main

import (
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/hnasser-dev/wordle-solver/internal/game"
	"github.com/hnasser-dev/wordle-solver/internal/words"
)

// Simulate a game for each possible answer
func main() {

	// file logger
	logFilePath := "logs/simulate.log"
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("unable to open log file %s - err: %s", logFilePath, err)
	}
	logger := slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	wordList, err := words.GetWordList()
	if err != nil {
		slog.Error("unable to read word list", "err", err)
		os.Exit(1)
	}

	// simulate a game for all possible answers
	var wg sync.WaitGroup
	for idx, answer := range wordList {
		wg.Add(1)
		go func() {
			defer wg.Done()
			guesses, gameWon := game.PlayGame(answer, wordList)
			slog.Info(
				"game complete",
				slog.Int("gameNum", idx+1),
				slog.String("answer", answer),
				slog.Int("numGuesses", len(guesses)),
				slog.Bool("gameWon", gameWon),
			)
		}()
	}
	wg.Wait()
}
