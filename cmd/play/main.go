package main

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"sort"
	"strings"

	"github.com/messy-coding/wordle-solver/internal/words"
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

var correctGuessColourPattern = colourPattern{Green, Green, Green, Green, Green}

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

func getSortedGuessOutcomes(remainingWords []string) []guessOutcome {

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

	sort.Slice(guessOutcomes, func(i, j int) bool { return guessOutcomes[i].entropyBits > guessOutcomes[j].entropyBits })

	return guessOutcomes
}

func playGame(answer string, wordList []string) error {

	remainingWordList := make([]string, len(wordList))
	copy(remainingWordList, wordList)

	guesses := []string{}
	gameWon := false

	for !gameWon && len(guesses) < maxNumGuesses {

		log.Printf("== guess #%d ==", len(guesses)+1)

		log.Printf("len(remainingWordList)=%d", len(remainingWordList))
		bestOutcomes := getSortedGuessOutcomes(remainingWordList)

		// just for logging
		topN := int(math.Min(3, float64(len(bestOutcomes))))
		topNAsStr := make([]string, topN)
		for i := range topN {
			outcome := bestOutcomes[i]
			topNAsStr[i] = fmt.Sprintf("guessOutcome(guess=%s, entropyBits=%f)", outcome.guess, outcome.entropyBits)
		}
		log.Printf("top %d next outcomes: %v", topN, strings.Join(topNAsStr, ","))

		bestOutcome := bestOutcomes[0]
		guesses = append(guesses, bestOutcome.guess)
		log.Printf("guessing word: %q", bestOutcome.guess)

		nextColourPattern := getColourPattern(bestOutcome.guess, answer)
		if reflect.DeepEqual(nextColourPattern, correctGuessColourPattern) {
			gameWon = true
			log.Printf("%q is the correct answer!", bestOutcome.guess)
			break
		}

		remainingWordList = bestOutcome.distribution[nextColourPattern]
	}

	log.Printf("guesses: %v", guesses)

	if gameWon {
		log.Printf("You won in %d guesses!", len(guesses))
	} else {
		log.Println("You lost!")
	}

	return nil

}

func main() {
	wordList, err := words.GetWordList()
	if err != nil {
		log.Fatal(err)
	}

	answer := "hello"
	log.Printf("(The answer is: %q - shhhhhh!)", answer)

	if err := playGame(answer, wordList); err != nil {
		log.Fatal(err)
	}
}
