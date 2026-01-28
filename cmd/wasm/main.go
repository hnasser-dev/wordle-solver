package main

import (
	"fmt"
	"syscall/js"

	"github.com/hnasser-dev/wordle-solver/internal/game"
	"github.com/hnasser-dev/wordle-solver/internal/words"
)

var wordList []string
var freqMap words.WordFrequencyMap

type nyTimesApiResponse struct {
	Id              int    `json:"id"`
	Solution        string `json:"solution"`
	PrintDate       string `json:"print_date"`
	DaysSinceLaunch int    `json:"days_since_launch"`
	Editor          string `json:"editor"`
}

// solveWordle(mode, initialGuesses...) -> {gameWon bool, guesses []string}
func solveWordle(this js.Value, args []js.Value) interface{} {
	answer := args[0].String()
	gameModeStr := args[1].String()
	var gameMode game.GameMode
	switch gameModeStr {
	case "normal":
		gameMode = game.NormalMode
	case "dumb":
		gameMode = game.DumbMode
	default:
		return js.ValueOf(fmt.Sprintf("error: invalid game mode type: %s"))
	}
	initialGuesses := []string{}
	if len(args) > 2 {
		for _, arg := range args[2:] {
			initialGuesses = append(initialGuesses, arg.String())
		}
	}
	game, err := game.NewGame(
		game.GameConfig{
			Answer:         answer,
			GameMode:       gameMode,
			InitialGuesses: initialGuesses,
			WordList:       wordList,
			FreqMap:        freqMap,
		})
	if err != nil {
		return js.ValueOf(fmt.Sprintf("error: unable to create wordle game - err: %s", err))
	}
	gameWon, guesses := game.PlayGameUntilEnd(true)
	jsObj := js.Global().Get("Object").New()
	jsObj.Set("gameWon", gameWon)
	jsGuesses := js.Global().Get("Array").New(len(guesses))
	for i, guess := range guesses {
		jsGuesses.SetIndex(i, guess)
	}
	jsObj.Set("guesses", jsGuesses)
	return jsObj
}

func main() {
	var err error
	wordList = words.GetWordList()
	freqMap, err = words.GetWordFrequencyMap()
	if err != nil {
		js.Global().Set("freqMapError", js.ValueOf(fmt.Sprintf("unable to fetch word frequency map - err: %s", err)))
	}
	js.Global().Set("solveWordle", js.FuncOf(solveWordle))
	select {}
}
