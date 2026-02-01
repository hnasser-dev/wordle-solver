//go:build ignore

package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/hnasser-dev/wordle-solver/internal/game"
)

func main() {
	guessHelper, err := game.NewGuessHelper(game.GuessHelperConfig{})
	if err != nil {
		panic(err)
	}
	sortedFirstGuessOutcomes := guessHelper.GetSortedGuessOutcomes(game.NormalMode)
	sortedOptimalGuesses := make([]string, len(sortedFirstGuessOutcomes))
	for i, outcome := range sortedFirstGuessOutcomes {
		sortedOptimalGuesses[i] = outcome.Guess
	}
	outputPath := "data/optimal_first_guesses.txt"
	err = os.WriteFile(outputPath, []byte(strings.Join(sortedOptimalGuesses, "\n")), 0644)
	if err != nil {
		panic(err)
	}
	slog.Info("written sortedOptimalGuesses", "path", outputPath)
}
