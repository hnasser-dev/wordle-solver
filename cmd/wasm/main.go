package main

import (
	"fmt"
	"syscall/js"

	"github.com/hnasser-dev/wordle-solver/internal/game"
	"github.com/hnasser-dev/wordle-solver/internal/words"
)

var wordList []string
var freqMap words.WordFrequencyMap

func sliceToJsArray[T any](s []T) js.Value {
	arr := make([]any, len(s))
	for i, v := range s {
		arr[i] = v
	}
	return js.ValueOf(arr)
}

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

	// default normal mode
	gameMode := game.NormalMode

	sortedOptimalFirstGuesses := words.GetOptimalFirstGuessesList()
	topN := 100
	jsOptimalFirstGuesses := js.Global().Get("Array").New()
	for idx, guess := range sortedOptimalFirstGuesses {
		if idx >= topN {
			break
		}
		jsOptimalFirstGuesses.Call("push", guess)
	}

	js.Global().Set("optimalFirstGuesses", jsOptimalFirstGuesses)

	jsGuessHelper := js.Global().Get("Object").New()
	// getSuggestions(guess string, colourPattern []string)
	getSuggestions := js.FuncOf(func(_ js.Value, args []js.Value) any {
		if len(args) != 2 {
			return js.Global().Get("Error").New("Incorrect number of arguments to getSuggestions - must be 2")
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
		guessHelper.FilterRemainingWords(guess, colourPattern)
		sortedGuessOutcomes := guessHelper.GetSortedGuessOutcomes(gameMode)
		returnArr := js.Global().Get("Array").New()
		for _, guessOutcome := range sortedGuessOutcomes {
			returnArr.Call("push", guessOutcome.Guess)
		}
		return returnArr
	})

	jsGuessHelper.Set("getSuggestions", getSuggestions)
	js.Global().Set("guessHelper", jsGuessHelper)

	select {}
}
