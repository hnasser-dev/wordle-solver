package game

import (
	"fmt"
	"reflect"

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

func NewGame(answer string) (*Game, error) {

	initialWordList, err := words.GetWordList()
	if err != nil {
		return nil, err
	}
	remainingWordList := make([]string, len(initialWordList))
	copy(remainingWordList, initialWordList)

	answerInWordList := false
	for _, word := range initialWordList {
		if word == answer {
			answerInWordList = true
			break
		}
	}

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
func (g *Game) Guess(guess string) bool {
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
}
