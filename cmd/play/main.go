package main

import (
	"log"

	"github.com/hnasser-dev/wordle-solver/internal/game"
	"github.com/hnasser-dev/wordle-solver/internal/words"
)

func main() {
	wordList, err := words.GetWordList()
	if err != nil {
		log.Fatal(err)
	}

	answer := "hello"
	log.Printf("(The answer is: %q - shhhhhh!)", answer)

	if err := game.PlayGame(answer, wordList); err != nil {
		log.Fatal(err)
	}
}
