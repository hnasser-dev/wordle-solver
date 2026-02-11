/*
Struct for suggesting optimal (or suboptimal) guesses, operating only off the pattern observed
The game's answer is NOT known ahead of time
*/

package game

import (
	"errors"
	"fmt"
	"maps"

	"github.com/hnasser-dev/wordle-solver/internal/words"
)

var ErrNoGuesses = errors.New("no guesses have been made")

type GuessHelperConfig struct {
	AllPossibleAnswers []string
	FreqMap            words.WordFrequencyMap
}

type GuessHelper struct {
	// GameMode                 GameMode
	FreqMap                     words.WordFrequencyMap
	AllRemainingPossibleAnswers [][]string // length: len(Guesses) + 1
	AllGuesses                  []string
	AllSortedGuessOutcomes      [][]guessOutcome
}

func NewGuessHelper(config GuessHelperConfig) (*GuessHelper, error) {
	var err error
	var allPossibleAnswers []string
	if config.AllPossibleAnswers == nil {
		allPossibleAnswers = words.GetPossibleAnswers()
	} else {
		allPossibleAnswers = append([]string{}, config.AllPossibleAnswers...) // copy
	}
	freqMap := words.WordFrequencyMap{}
	if config.FreqMap == nil {
		freqMap, err = words.GetWordFrequencyMap()
		if err != nil {
			return nil, fmt.Errorf("unable to read frequency map - err: %w", err)
		}
	} else {
		maps.Copy(freqMap, config.FreqMap)
	}
	guessHelper := GuessHelper{FreqMap: freqMap, AllRemainingPossibleAnswers: [][]string{allPossibleAnswers}}
	return &guessHelper, nil
}

func (g *GuessHelper) MakeGuess(guess string, pattern colourPattern) {
	possibleAnswers := g.AllRemainingPossibleAnswers[len(g.AllRemainingPossibleAnswers)-1]
	guessDistribution := computeGuessDistribution(guess, possibleAnswers)
	nextPossibleAnswers := guessDistribution[pattern]
	sortedGuessOutcomes := GetSortedGuessOutcomes(nextPossibleAnswers, g.FreqMap)
	g.AllRemainingPossibleAnswers = append(g.AllRemainingPossibleAnswers, nextPossibleAnswers)
	g.AllGuesses = append(g.AllGuesses, guess)
	g.AllSortedGuessOutcomes = append(g.AllSortedGuessOutcomes, sortedGuessOutcomes)
}

func (g *GuessHelper) RevertLastGuess() error {
	if len(g.AllGuesses) == 0 {
		return ErrNoGuesses
	}
	g.AllGuesses = g.AllGuesses[:len(g.AllGuesses)-1]
	g.AllRemainingPossibleAnswers = g.AllRemainingPossibleAnswers[:len(g.AllRemainingPossibleAnswers)-1]
	g.AllSortedGuessOutcomes = g.AllSortedGuessOutcomes[:len(g.AllSortedGuessOutcomes)-1]
	return nil
}
