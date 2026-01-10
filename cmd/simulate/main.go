package main

import (
	"log/slog"
	"os"
	"sync"

	"github.com/hnasser-dev/wordle-solver/internal/game"
	"github.com/hnasser-dev/wordle-solver/internal/words"
)

// Simulate a game for each possible answer
func main() {
	logger := slog.New(slog.NewJSONHandler(os.NewFile()))

	wordList, err := words.GetWordList()
	if err != nil {
		slog.Error("unable to read word list", "err", err)
		os.Exit(1)
	}
	var wg sync.WaitGroup
	for idx, answer := range wordList[:50] {
		wg.Add(1)
		go func() {
			defer wg.Done()
			guesses, gameWon := game.PlayGame(answer, wordList)
			slog.Info("game complete", "num", idx+1, "answer", answer, "numGuesses", len(guesses), "gameWon", gameWon)
		}()
	}
	wg.Wait()
}
