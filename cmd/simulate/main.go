package main

import (
	"log"
	"log/slog"
	"os"
	"runtime"
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
	fileLogger := slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slog.LevelInfo}))

	wordList, err := words.GetWordList()
	if err != nil {
		slog.Error("unable to read word list", "err", err)
		os.Exit(1)
	}

	// simulate a game for all possible answers - multiprocess and limit concurrency
	var wg sync.WaitGroup
	sem := make(chan struct{}, runtime.NumCPU())

	for idx, answer := range wordList {
		wg.Add(1)
		sem <- struct{}{} // block until there is space in the sem
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			guesses, gameWon := game.PlayGame(answer, wordList)
			fileLogger.Info(
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
