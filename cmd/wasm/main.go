package main

import (
	"fmt"
	"log"
	"syscall/js"

	"github.com/hnasser-dev/wordle-solver/internal/game"
	"github.com/hnasser-dev/wordle-solver/internal/words"
)

var wordList []string
var freqMap words.WordFrequencyMap

func main() {

	var err error

	wordList = words.GetWordList()

	freqMap, err = words.GetWordFrequencyMap()
	if err != nil {
		js.Global().Get("Error").New(fmt.Sprintf("unable to fetch word frequency map - err: %s", err))
	}

	guessHelper, err := game.NewGuessHelper(game.GuessHelperConfig{WordList: wordList, FreqMap: freqMap})
	if err != nil {
		js.Global().Get("Error").New(fmt.Sprintf("unable to create guessHelper - err: %s", err))
	}

	obj := js.Global().Get("Object").New()

	// getSuggestions(guess string, colourPattern []string, gameMode string)
	getSuggestions := js.FuncOf(func(_ js.Value, args []js.Value) any {
		log.Println("hi!!")
		if len(args) != 3 {
			return js.Global().Get("Error").New("Incorrect number of arguments to getSuggestions - must be 3")
		}
		guess := args[0].String()
		jsColourStringsArr := args[1]
		if !jsColourStringsArr.InstanceOf(js.Global().Get("Array")) {
			return js.Global().Get("Error").New("args[1] must be an array")
		}
		colourStringsLength := jsColourStringsArr.Length()
		if colourStringsLength != game.WordLength {
			return js.Global().Get("Error").New(fmt.Sprintf("colour pattern must be of length %d"))
		}
		var colourStringsSlice [game.WordLength]string
		for i := 0; i < len(colourStringsSlice); i++ {
			colourStringsSlice[i] = jsColourStringsArr.Index(i).String()
		}
		colourPattern, err := game.ColourStringsToColourPattern(colourStringsSlice)
		if err != nil {
			return js.Global().Get("Error").New(fmt.Sprintf("unable to parse colour strings: %s", err))
		}
		gameMode, err := game.GameModeStringToGameMode(args[2].String())
		if err != nil {
			return js.Global().Get("Error").New(fmt.Sprintf("unable to parse game mode: %s", err))
		}
		guessHelper.FilterRemainingWords(guess, colourPattern)
		return guessHelper.GetSortedGuessOutcomes(gameMode)
	})

	obj.Set("getSuggestions", getSuggestions)
	js.Global().Set("guessHelper", obj)

	select {}
}
