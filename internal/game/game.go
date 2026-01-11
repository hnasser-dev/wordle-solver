package game

import (
	"fmt"
	"log/slog"
	"math"
	"reflect"
	"sort"
	"strings"

	"github.com/hnasser-dev/wordle-solver/internal/words"
)

const wordLength int = 5

const maxNumGuesses int = 6

const (
	Grey colour = iota
	Yellow
	Green
)

type colour uint8

type colourPattern [wordLength]colour

type guessDistribution map[colourPattern][]string

type guessOutcome struct {
	guess        string
	distribution guessDistribution
	entropyBits  float64
}

type gameState struct {
	guesses               []string
	bestRemainingOutcomes []guessOutcome
	gameWon               bool
}

var correctGuessColourPattern = colourPattern{Green, Green, Green, Green, Green}

// Returns the guesses and whether the game was won
func PlayGame(answer string, wordList []string, freqMap words.WordFrequencyMap, initialGameState *gameState) ([]string, bool) {

	remainingWordList := make([]string, len(wordList))
	copy(remainingWordList, wordList)

	state := initialGameState
	if state == nil {
		state = &gameState{
			guesses:               []string{},
			gameWon:               false,
			bestRemainingOutcomes: []guessOutcome{},
		}
	}

	for !state.gameWon && len(state.guesses) < maxNumGuesses {

		slog.Debug("guess", slog.Int("number", len(state.guesses)+1))

		slog.Debug("len(remainingWordList)", slog.Int("len", len(remainingWordList)))
		bestOutcomes := getSortedGuessOutcomes(remainingWordList, freqMap)

		// just for logging
		topN := int(math.Min(3, float64(len(bestOutcomes))))
		topNAsStr := make([]string, topN)
		for i := range topN {
			outcome := bestOutcomes[i]
			topNAsStr[i] = fmt.Sprintf("guessOutcome(guess=%s, entropyBits=%f)", outcome.guess, outcome.entropyBits)
		}
		slog.Debug(
			"top next outcomes",
			slog.Int("topN", topN),
			slog.String("outcomes", strings.Join(topNAsStr, ",")),
		)

		bestOutcome := bestOutcomes[0]
		state.guesses = append(state.guesses, bestOutcome.guess)
		slog.Debug("performing next guess", slog.String("guess", bestOutcome.guess))

		nextColourPattern := getColourPattern(bestOutcome.guess, answer)
		if reflect.DeepEqual(nextColourPattern, correctGuessColourPattern) {
			state.gameWon = true
			slog.Debug("correct answer", slog.String("answer", bestOutcome.guess))
			break
		}

		remainingWordList = bestOutcome.distribution[nextColourPattern]
	}

	return state.guesses, state.gameWon
}

func getColourPattern(guess string, answer string) colourPattern {

	colourPatternSlice := [wordLength]colour{Grey, Grey, Grey, Grey, Grey}

	numCharsInAnswer := map[byte]int{}
	for i := 0; i < len(answer); i++ {
		numCharsInAnswer[answer[i]]++
	}

	numCharsGuessed := map[byte]int{}

	// assign greens
	for idx := range wordLength {
		if guess[idx] == answer[idx] {
			colourPatternSlice[idx] = Green
			numCharsGuessed[guess[idx]]++
		}
	}

	// assign yellows
	for idx := range wordLength {
		if colourPatternSlice[idx] == Green {
			continue
		}
		guessChar := guess[idx]
		if numCharsInAnswer[guessChar] > numCharsGuessed[guessChar] {
			colourPatternSlice[idx] = Yellow
			numCharsGuessed[guessChar]++
		}
	}

	return colourPatternSlice
}

func getSortedGuessOutcomes(remainingWords []string, freqMap words.WordFrequencyMap) []guessOutcome {

	guessDistributions := map[string]guessDistribution{}
	for _, potentialGuess := range remainingWords {
		dist := guessDistribution{}
		for _, potentialAnswer := range remainingWords {
			if potentialAnswer == potentialGuess {
				continue // prevent infinite loop of guesses
			}
			colourPattern := getColourPattern(potentialGuess, potentialAnswer)
			dist[colourPattern] = append(dist[colourPattern], potentialAnswer)
		}
		guessDistributions[potentialGuess] = dist
	}

	guessOutcomes := make([]guessOutcome, 0, len(guessDistributions))
	for guess, dist := range guessDistributions {
		outcome := guessOutcome{
			guess:        guess,
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
				return freqMap[guessOutcomes[i].guess] > freqMap[guessOutcomes[j].guess]
			} else {
				return guessOutcomes[i].entropyBits > guessOutcomes[j].entropyBits
			}
		},
	)

	return guessOutcomes
}
