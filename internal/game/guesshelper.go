/*
Struct for suggesting optimal (or suboptimal) guesses, operating only off the pattern observed
The game's answer is NOT known ahead of time
*/

package game

import (
	"errors"
	"fmt"
	"maps"
	"slices"

	"github.com/hnasser-dev/wordle-solver/internal/words"
)

var ErrNoGuesses = errors.New("no guesses have been made")

type GuessHelperConfig struct {
	WordList []string
	FreqMap  words.WordFrequencyMap
}

type GuessHelper struct {
	FreqMap   words.WordFrequencyMap
	Guesses   []string
	WordLists [][]string // length: len(Guesses) + 1
}

func NewGuessHelper(config GuessHelperConfig) (*GuessHelper, error) {
	var err error
	var initialWordList []string
	if config.WordList == nil {
		initialWordList = words.GetPossibleAnswers()
	} else {
		initialWordList = append([]string{}, config.WordList...) // copy
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
	guessHelper := GuessHelper{FreqMap: freqMap, WordLists: [][]string{initialWordList}}
	return &guessHelper, nil
}

func (g *GuessHelper) MakeGuess(guess string, pattern colourPattern) {
	remainingWords := g.WordLists[len(g.Guesses)]
	guessDistribution := computeGuessDistribution(guess, remainingWords)
	g.WordLists = append(g.WordLists, guessDistribution[pattern])
	g.Guesses = append(g.Guesses, guess)
}

func (g *GuessHelper) GetSortedGuessOutcomes(gameMode GameMode) []guessOutcome {
	remainingWords := g.WordLists[len(g.Guesses)]
	sortedGuessOutcomes := getSortedGuessOutcomes(remainingWords, g.FreqMap)
	switch gameMode {
	case DumbMode:
		slices.Reverse(sortedGuessOutcomes)
	case NormalMode:
		// do nothing
	default:
		panic(fmt.Sprintf("unknown gameMode: %d", gameMode))
	}
	return sortedGuessOutcomes
}

func (g *GuessHelper) RevertLastGuess() error {
	if len(g.Guesses) == 0 {
		return ErrNoGuesses
	}
	g.Guesses = g.Guesses[:len(g.Guesses)-1]
	g.WordLists = g.WordLists[:len(g.WordLists)-1]
	return nil
}
