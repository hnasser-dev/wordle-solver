package main

import (
	"log"
	"math"
	"reflect"
	"sort"

	"github.com/messy-coding/wordle-solver/internal/words"
)

const wordLength int = 5

const maxNumGuesses int = 6

type Colour uint8

const (
	Grey Colour = iota
	Yellow
	Green
)

type guessEntropy struct {
	Guess string
	Bits  float64
}

var correctGuessColourPattern = [wordLength]Colour{Green, Green, Green, Green, Green}

func getColourPattern(guess string, answer string) [wordLength]Colour {

	colourPatternSlice := [wordLength]Colour{Grey, Grey, Grey, Grey, Grey}

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

// TODO - fix this so it works - currently gets stuck - and not sure if this even makes sense anyway
func filterWordListBasedOnLatestGuess(wordList []string, guessColourPattern [wordLength]Colour, prevGuess string, answer string) []string {
	filteredWordList := []string{}
	for _, word := range wordList {
		// prevent repeating the same guess
		if word == prevGuess {
			continue
		}
		wordColourPattern := getColourPattern(word, answer)
		// log.Printf("%s is %v", word, wordColourPattern)
		wordCouldBeAnswer := reflect.DeepEqual(guessColourPattern, wordColourPattern)
		if wordCouldBeAnswer {
			// log.Printf("%s + %v could be %s", word, guessColourPattern, answer)
			filteredWordList = append(filteredWordList, word)
		}
	}
	return filteredWordList
}

func getOptimalNextGuess(remainingWords []string) string {

	entropies := map[string]float64{}

	// TODO - fix entropies calculation: I think it may be wrong
	for _, potentialGuess := range remainingWords {

		colourPatternNums := map[[wordLength]Colour]int{}

		for _, potentialAnswer := range remainingWords {
			colourPattern := getColourPattern(potentialGuess, potentialAnswer)
			colourPatternNums[colourPattern]++
		}

		numRemainingWords := len(remainingWords)
		numBitsOfInfo := 0.0

		for _, numWordsInPath := range colourPatternNums {
			probability := float64(numWordsInPath) / float64(numRemainingWords)
			numBitsOfInfo += float64(probability) * math.Log2(1.0/probability)
		}
		entropies[potentialGuess] = numBitsOfInfo
	}

	entropiesHighToLow := make([]guessEntropy, 0, len(entropies))

	for guess, bits := range entropies {
		entropiesHighToLow = append(entropiesHighToLow, guessEntropy{guess, bits})
	}

	sort.Slice(entropiesHighToLow, func(i, j int) bool { return entropiesHighToLow[i].Bits > entropiesHighToLow[j].Bits })

	// TODO - remove
	nextGuess := entropiesHighToLow[0].Guess
	log.Printf("Next guess: %s (Top 10 pairs: %v)", nextGuess, entropiesHighToLow[:10])

	entropiesLowToHigh := make([]guessEntropy, 0, len(entropies))
	entropiesLowToHigh = append(entropiesLowToHigh, entropiesHighToLow...)
	sort.Slice(entropiesLowToHigh, func(i, j int) bool { return entropiesLowToHigh[i].Bits < entropiesLowToHigh[j].Bits })
	log.Printf("Bottom 10 pairs: %v", entropiesLowToHigh[:10])

	return nextGuess
}

func playGame(answer string, wordList []string) error {

	remainingWordList := make([]string, len(wordList))
	copy(remainingWordList, wordList)

	numGuesses := 0
	gameWon := false

	for numGuesses < maxNumGuesses {
		numGuesses++
		guess := getOptimalNextGuess(remainingWordList)
		guessColourPattern := getColourPattern(guess, answer)
		if reflect.DeepEqual(guessColourPattern, correctGuessColourPattern) {
			gameWon = true
			break
		}
		remainingWordList = filterWordListBasedOnLatestGuess(remainingWordList, guessColourPattern, guess, answer)
		log.Printf("len(remainingWordList)=%d", len(remainingWordList))
	}

	if gameWon {
		log.Printf("You won in %d guesses!", numGuesses)
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

	log.Println(getOptimalNextGuess(wordList))

	answer := "hello"
	log.Printf("(The answer is: %q - shhhhhh!)", answer)

	if err := playGame(answer, wordList); err != nil {
		log.Fatal(err)
	}
}
