package main

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"sort"

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

func (g guessOutcome) String() string {
	return fmt.Sprintf("guessOutcome(guess=%s, entropyBits=%f)", g.guess, g.entropyBits)
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

	guessDistributions := make(map[string]guessDistribution, len(remainingWords))

	for _, potentialGuess := range remainingWords {
		dist := guessDistribution{}
		for _, potentialAnswer := range remainingWords {
			colourPattern := getColourPattern(potentialGuess, potentialAnswer)
			dist[colourPattern] = append(dist[colourPattern], potentialAnswer)
		}
		guessDistributions[potentialGuess] = dist
	}

	guessOutcomes := make([]guessOutcome, len(guessDistributions))

	// entropy calculations
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

	// sort
	sort.Slice(guessOutcomes, func(i, j int) bool { return guessOutcomes[i].entropyBits > guessOutcomes[j].entropyBits })
	return guessOutcomes
}

func playGame(answer string, wordList []string) error {

	remainingWordList := make([]string, len(wordList))
	copy(remainingWordList, wordList)

	guesses := []string{}
	gameWon := false

	for !gameWon && len(guesses) < maxNumGuesses {

		nextOutcomes := getSortedGuessOutcomes(wordList)
		for _, outcome := range nextOutcomes[:3] {
			log.Println(outcome)
		}

		nextOutcome := nextOutcomes[0]
		guesses = append(guesses, nextOutcome.guess)

		nextColourPattern := getColourPattern(nextOutcome.guess, answer)
		if reflect.DeepEqual(nextColourPattern, correctGuessColourPattern) {
			gameWon = true
			break
		}

		remainingWordList = nextOutcome.distribution[nextColourPattern]
		log.Printf("len(remainingWordList)=%d", len(remainingWordList))
		log.Printf("guesses: %v", guesses)
	}

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
