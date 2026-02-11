package game

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/hnasser-dev/wordle-solver/internal/words"
)

const (
	NormalMode GameMode = iota
	DumbMode
)

const (
	Grey colour = iota
	Yellow
	Green
)

type colour uint8
type colourPattern [WordLength]colour

type guessDistribution map[colourPattern][]string

type guessOutcome struct {
	Guess        string
	distribution guessDistribution
	entropyBits  float64
}

type GameMode uint8

func (m GameMode) Valid() bool {
	switch m {
	case NormalMode, DumbMode:
		return true
	default:
		return false
	}
}

func GetSortedGuessOutcomes(remainingWords []string, freqMap words.WordFrequencyMap) []guessOutcome {

	guessDistributions := map[string]guessDistribution{}
	for _, potentialGuess := range remainingWords {
		guessDistributions[potentialGuess] = computeGuessDistribution(potentialGuess, remainingWords)
	}

	guessOutcomes := make([]guessOutcome, 0, len(guessDistributions))
	for guess, dist := range guessDistributions {
		outcome := guessOutcome{
			Guess:        guess,
			distribution: dist,
			entropyBits:  0.0,
		}
		for _, guesses := range dist {
			prob := float64(len(guesses)) / float64(len(remainingWords))
			outcome.entropyBits += float64(prob) * math.Log2(1.0/prob)
		}
		guessOutcomes = append(guessOutcomes, outcome)
	}

	sort.Slice(
		guessOutcomes,
		func(i, j int) bool {
			// if equal entropies, prioritise higher frequency
			if guessOutcomes[i].entropyBits == guessOutcomes[j].entropyBits {
				return freqMap[guessOutcomes[i].Guess] > freqMap[guessOutcomes[j].Guess]
			} else {
				return guessOutcomes[i].entropyBits > guessOutcomes[j].entropyBits
			}
		},
	)

	return guessOutcomes
}

func ColourStringsToColourPattern(colourSlice [WordLength]string) (colourPattern, error) {
	var c colourPattern
	for i, colourString := range colourSlice {
		colourString = strings.ToLower(colourString)
		switch colourString {
		case "green":
			c[i] = Green
		case "yellow":
			c[i] = Yellow
		case "grey":
			c[i] = Grey
		default:
			return c, fmt.Errorf("unknown colour: %s", colourString)
		}
	}
	return c, nil
}

func computeGuessDistribution(guess string, wordList []string) guessDistribution {
	dist := guessDistribution{}
	for _, potentialAnswer := range wordList {
		if potentialAnswer == guess {
			continue
		}
		colourPattern := getColourPattern(guess, potentialAnswer)
		dist[colourPattern] = append(dist[colourPattern], potentialAnswer)
	}
	return dist
}
