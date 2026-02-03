package main

import (
	"log"
	"log/slog"
	"os"
	"runtime"
	"sort"
	"sync"

	"github.com/hnasser-dev/wordle-solver/internal/game"
	"github.com/hnasser-dev/wordle-solver/internal/words"
)

type avgGuessesEntry struct {
	word       string
	avgGuesses float64
}

type winRatesEntry struct {
	word    string
	winRate float64
}

// Simulate a game for each possible answer
func main() {

	// file logger
	logFilePath := "logs/simulate.log"
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("unable to open log file %s - err: %s", logFilePath, err)
	}
	fileLogger := slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slog.LevelInfo}))

	possibleAnswers := words.GetPossibleAnswers()

	freq, err := words.GetWordFrequencyMap()
	if err != nil {
		slog.Error("unable to read frequency map list", "err", err)
		os.Exit(1)
	}

	// independent subset of list
	startingGuessesList := append([]string{}, possibleAnswers[:8]...)

	allAvgGuesses := make([]avgGuessesEntry, 0, len(startingGuessesList))
	allWinRates := make([]winRatesEntry, 0, len(startingGuessesList))

	// simulate a game for all possible startingGuesses and answers
	var wg sync.WaitGroup
	sem := make(chan struct{}, runtime.NumCPU())

	slog.Info("starting simulation", slog.Int("numStartingGuesses", len(startingGuessesList)))
	for _, startingGuess := range startingGuessesList {
		totalNumGames := 0
		totalNumGuesses := 0
		totalNumWins := 0
		for _, answer := range possibleAnswers {
			wg.Add(1)
			sem <- struct{}{} // block until there is space in the sem
			go func() {
				defer wg.Done()
				defer func() { <-sem }()
				initialGuesses := []string{startingGuess}
				game, err := game.NewGameSimulator(game.GameSimulatorConfig{
					Answer:         answer,
					GameMode:       game.NormalMode,
					InitialGuesses: initialGuesses,
					WordList:       possibleAnswers,
					FreqMap:        freq,
				})
				if err != nil {
					slog.Error("unable to create new game", "err", err)
					os.Exit(1)
				}
				gameWon, guesses := game.PlayGameUntilEnd(false)
				totalNumGames++
				totalNumGuesses += len(guesses)
				if gameWon {
					totalNumWins++
				}
			}()
		}
		avgGuesses := float64(totalNumGuesses) / float64(totalNumWins)
		allAvgGuesses = append(
			allAvgGuesses,
			avgGuessesEntry{
				word:       startingGuess,
				avgGuesses: avgGuesses,
			},
		)
		winRate := float64(totalNumWins) / float64(totalNumGames)
		allWinRates = append(
			allWinRates,
			winRatesEntry{
				word:    startingGuess,
				winRate: winRate,
			},
		)
		fileLogger.Info(
			"iteration complete",
			slog.String("startingGuess", startingGuess),
			slog.Float64("avgGuess", avgGuesses),
			slog.Float64("winRate", winRate),
		)
	}
	slog.Info("completed simulation", slog.Int("numStartingGuesses", len(startingGuessesList)))

	// sorted allAvgGuesses, high to low
	sort.Slice(
		allAvgGuesses,
		func(i, j int) bool { return allAvgGuesses[i].avgGuesses < allAvgGuesses[j].avgGuesses },
	)

	// sorted allWinRates, high to low
	sort.Slice(
		allWinRates,
		func(i, j int) bool { return allWinRates[i].winRate > allWinRates[j].winRate },
	)

	slog.Info("best words by avgGuesses", "wordsAndAvgGuesses", allAvgGuesses[:5])
	slog.Info("best words by winRate", "wordsAndWinRates", allWinRates[:5])

	wg.Wait()
}
