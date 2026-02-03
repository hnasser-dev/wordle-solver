/*
Struct for suggesting optimal (or suboptimal) guesses, operating only off the pattern observed
The game's answer is NOT known ahead of time
*/

package game

import (
	"fmt"
	"maps"
	"slices"

	"github.com/hnasser-dev/wordle-solver/internal/words"
)

type GuessHelperConfig struct {
	WordList []string
	FreqMap  words.WordFrequencyMap
}

type GuessHelper struct {
	FreqMap        words.WordFrequencyMap
	RemainingWords []string
}

func NewGuessHelper(config GuessHelperConfig) (*GuessHelper, error) {
	var err error
	var remainingWords []string
	if config.WordList == nil {
		remainingWords = words.GetPossibleAnswers()
	} else {
		remainingWords = append([]string{}, config.WordList...) // copy
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
	guessHelper := GuessHelper{FreqMap: freqMap, RemainingWords: remainingWords}
	return &guessHelper, nil
}

func (g *GuessHelper) FilterRemainingWords(guess string, pattern colourPattern) {
	guessDistribution := computeGuessDistribution(guess, g.RemainingWords)
	g.RemainingWords = guessDistribution[pattern]
}

func (g *GuessHelper) GetSortedGuessOutcomes(gameMode GameMode) []guessOutcome {
	sortedGuessOutcomes := getSortedGuessOutcomes(g.RemainingWords, g.FreqMap)
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
