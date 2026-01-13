package game

import (
	"fmt"
	"reflect"

	"slices"

	"github.com/hnasser-dev/wordle-solver/internal/words"
)

type Game struct {
	Answer                  string
	GameWon                 bool
	Guesses                 []string
	InitialWordList         []string
	RemainingWordList       []string
	SortedRemainingOutcomes []guessOutcome
	WordFrequencies         words.WordFrequencyMap
}

func NewGame(answer string, initialGuesses ...string) (*Game, error) {

	initialWordList, err := words.GetWordList()
	if err != nil {
		return nil, err
	}
	remainingWordList := make([]string, len(initialWordList))
	copy(remainingWordList, initialWordList)

	answerInWordList := slices.Contains(initialWordList, answer)
	if !answerInWordList {
		return nil, fmt.Errorf("provided answer %q is not in the word list", answer)
	}

	game := Game{
		Answer:                  answer,
		GameWon:                 false,
		Guesses:                 []string{},
		InitialWordList:         initialWordList,
		RemainingWordList:       remainingWordList,
		SortedRemainingOutcomes: []guessOutcome{},
	}

	return &game, nil
}

// Returns: gameWon (bool)
func (g *Game) PerformGuess(guess string) bool {
	g.Guesses = append(g.Guesses, guess)
	nextColourPattern := getColourPattern(guess, g.Answer)
	if reflect.DeepEqual(nextColourPattern, correctGuessColourPattern) {
		g.GameWon = true
	}
	// TODO - create a function that returns guess outcomes as a (unsorted) map
	// THEN UPDATE THE BELOW FUNCTION CALL
	// And then index into the map like: g.RemainingWordList = g.SortedRemainingOutcomes[guess].distribution[nextColourGuess]
	g.SortedRemainingOutcomes = getSortedGuessOutcomes(g.RemainingWordList, g.WordFrequencies)
	g.RemainingWordList = 
	return g.GameWon
}

// Returns: gameWon (bool)
func (g *Game) PerformOptimalGuess() bool {
	g.SortedRemainingOutcomes = getSortedGuessOutcomes(g.RemainingWordList, g.WordFrequencies)
	bestOutcome := g.SortedRemainingOutcomes[0]
	g.Guesses = append(g.Guesses, bestOutcome.guess)
	nextColourPattern := getColourPattern(bestOutcome.guess, g.Answer)
	if reflect.DeepEqual(nextColourPattern, correctGuessColourPattern) {
		g.GameWon = true
	}
	return g.GameWon
}

func (g *Game) PlayGameUntilEnd(limitGuesses bool) {
	// play the game until complete
	for !g.GameWon {
		if limitGuesses && len(g.Guesses) == maxNumGuesses {
			break
		}
		g.PerformOptimalGuess()
	}
}
