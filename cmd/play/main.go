package main

import (
	"log"
	"reflect"

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

var correctGuessColourPattern = [wordLength]Colour{Green, Green, Green, Green, Green}

func getOptimalNextGuess(remainingWords []string) string {
	/*
		Iterate through all remaining words
		For each word, calculate the colour distribution pattern across all other words
		...
	*/

	// bitsPerGuess := map[string]float64{}

	for _, potentialGuess := range remainingWords {
		colourPatternNums := map[[5]Colour]int{}
		for _, potentialAnswer := range remainingWords {
			colourPattern := getColourPattern(potentialGuess, potentialAnswer)
			colourPatternNums[colourPattern]++
		}
		// TODO - calculate numBitsOfInfo (weighted sum of probabilities, then convert to bits using -log(p))
		// bitsPerGuess[potentialGuess] = numBitsOfInfo
	}

	return "hello"
}

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

	log.Println(colourPatternSlice)
	log.Println(numCharsGuessed)

	return colourPatternSlice
}

func playGame(answer string, wordList []string) error {

	remainingWordList := make([]string, len(wordList))
	copy(remainingWordList, wordList)

	numGuesses := 0
	gameWon := false

	for numGuesses < maxNumGuesses {
		numGuesses++
		nextGuess := getOptimalNextGuess(remainingWordList)
		colourPattern := getColourPattern(nextGuess, answer)
		if reflect.DeepEqual(colourPattern, correctGuessColourPattern) {
			gameWon = true
			break
		}
		// TODO - update remainingWordList to only include words compatible with colourPattern
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

	answer := "hello"
	log.Printf("(The answer is: %q - shhhhhh!)", answer)

	if err := playGame(answer, wordList); err != nil {
		log.Fatal(err)
	}
}
