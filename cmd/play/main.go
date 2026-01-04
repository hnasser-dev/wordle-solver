package main

import (
	"log"
	"math/rand"

	"github.com/messy-coding/wordle-solver/internal/words"
)

const wordLength int = 5
const maxNumGuesses int = 6

func selectRandomWord(wordList []string) string {
	return wordList[rand.Intn(len(wordList))]
}

func playGame(wordList []string, answer string) error {

	answerCharCounter := map[byte]int{}
	for _, val := range answer {
		answerCharCounter[byte(val)] += 1
	}

	correctChars := [wordLength]byte{}
	incorrectCharPositions := map[byte]map[int]struct{}{}
	invalidChars := map[byte]struct{}{}

	guess := ""
	numGuesses := 0
	possibleWordsRemaining := wordList[:]

	for numGuesses < maxNumGuesses {

		numGuesses++

		guess = selectRandomWord(possibleWordsRemaining)
		log.Printf("(guess #%d) num possible words remaining: %d -> guessing: %q", numGuesses, len(possibleWordsRemaining), guess)

		if guess == answer {
			break
		}

		// increase knowledge pool
		guessCharCounter := map[byte]int{}
		for i, rne := range guess {
			char := byte(rne)
			charOccurencesInAnswer, charExistsInAnswer := answerCharCounter[char]
			if char == answer[i] {
				correctChars[i] = char
			} else if charExistsInAnswer && guessCharCounter[char] <= charOccurencesInAnswer {
				if _, ok := incorrectCharPositions[char]; !ok {
					incorrectCharPositions[char] = map[int]struct{}{}
				}
				incorrectCharPositions[char][i] = struct{}{}
			} else {
				invalidChars[char] = struct{}{}
			}
			guessCharCounter[char] += 1
		}

		// update valid words left
		tmpPossibleWordsRemaining := []string{}
		for _, word := range possibleWordsRemaining {
			wordIsValidNextGuess := true
			for i, rne := range word {
				char := byte(rne)
				// don't guess a word with a character in a grey position
				if _, exists := invalidChars[char]; exists {
					wordIsValidNextGuess = false
					break
				}
				// don't guess a word with a character in a yellow position
				if incorrectCharPos, exists := incorrectCharPositions[char]; exists {
					for incorrectPosition := range incorrectCharPos {
						if incorrectPosition == i {
							wordIsValidNextGuess = false
							break
						}
					}
				}
				// don't guess a word that doesn't agree with a known green character
				if correctChars[i] != 0 && correctChars[i] != char {
					wordIsValidNextGuess = false
					break
				}

			}
			if wordIsValidNextGuess {
				tmpPossibleWordsRemaining = append(tmpPossibleWordsRemaining, word)
			}
		}
		possibleWordsRemaining = tmpPossibleWordsRemaining[:]
	}

	if guess != answer {
		log.Printf("You failed! Your closest guess after %d guesses was: %q. The answer was %q.", numGuesses, guess, answer)
	} else {
		log.Printf("Game completed! Correctly guessed %q after %d guesses", answer, numGuesses)
	}

	return nil
}

func main() {
	wordList, err := words.GetWordList()
	if err != nil {
		log.Fatal(err)
	}

	answer := selectRandomWord(wordList)
	log.Printf("(The answer is: %q - shhhhhh!)", answer)

	playGame(wordList, answer)
}
