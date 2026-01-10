package main

import (
	"log"

	"log/slog"

	"github.com/hnasser-dev/wordle-solver/internal/game"
	"github.com/hnasser-dev/wordle-solver/internal/words"
)

func main() {
	wordList, err := words.GetWordList()
	if err != nil {
		log.Fatal(err)
	}

	answer := "hello"
	slog.Debug("answer", "answer", answer)

	game.PlayGame(answer, wordList)
}
