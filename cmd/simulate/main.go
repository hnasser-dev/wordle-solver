package main

import (
	"io"
	"log"

	"github.com/hnasser-dev/wordle-solver/internal/game"
	"github.com/hnasser-dev/wordle-solver/internal/words"
)

func main() {
	wordList, err := words.GetWordList()
	if err != nil {
		log.Fatal(err)
	}

	// disable logging
	oldLogger := log.Writer()
	log.SetOutput(io.Discard)

	for _, word := range wordList {
		game.PlayGame(word, wordList)
	}

	// re-enable logging
	log.SetOutput(oldLogger)

}
