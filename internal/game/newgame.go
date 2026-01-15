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
	guessDistribution := computeGuessDistribution(guess, g.RemainingWordList)
	colourPattern := getColourPattern(guess, g.Answer)
	if reflect.DeepEqual(colourPattern, correctGuessColourPattern) {
		g.GameWon = true
	}
	g.RemainingWordList = guessDistribution[colourPattern]
	g.SortedRemainingOutcomes = getSortedGuessOutcomes(g.RemainingWordList, g.WordFrequencies)
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

func computeGuessDistribution(guess string, wordList []string) guessDistribution {
	dist := guessDistribution{}
	for _, potentialAnswer := range wordList {
		if potentialAnswer == guess {
			continue
		}
		colourPattern := getColourPattern(potentialAnswer, potentialAnswer)
		dist[colourPattern] = append(dist[colourPattern], potentialAnswer)
	}
	return dist
}
