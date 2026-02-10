package main

import (
	"fmt"
	"log"
	"syscall/js"

	"github.com/hnasser-dev/wordle-solver/internal/game"
	"github.com/hnasser-dev/wordle-solver/internal/words"
)

var (
	guessHelper *game.GuessHelper

	jsResetGuessHelper  js.Func
	jsGetSuggestedWords js.Func
	jsUndoLastGuess     js.Func
)

func sliceToJsArray[T any](s []T) js.Value {
	arr := make([]any, len(s))
	for i, v := range s {
		arr[i] = v
	}
	return js.ValueOf(arr)
}

func main() {

	var err error

	possibleAnswers := words.GetPossibleAnswers()

	freqMap, err := words.GetWordFrequencyMap()
	if err != nil {
		js.Global().Get("Error").New(fmt.Sprintf("unable to fetch word frequency map - err: %s", err))
	}

	gh, err := game.NewGuessHelper(game.GuessHelperConfig{WordList: possibleAnswers, FreqMap: freqMap})
	if err != nil {
		js.Global().Get("Error").New(fmt.Sprintf("unable to create guessHelper - err: %s", err))
	}
	guessHelper = gh

	validGuesses := words.GetValidGuesses()
	jsAllValidGuesses := js.Global().Get("Set").New()
	for _, guess := range validGuesses {
		jsAllValidGuesses.Call("add", guess)
	}

	sortedOptimalFirstGuesses := words.GetOptimalFirstGuessesList()
	topN := 100
	jsOptimalFirstGuesses := js.Global().Get("Array").New()
	for idx, guess := range sortedOptimalFirstGuesses {
		if idx >= topN {
			break
		}
		jsOptimalFirstGuesses.Call("push", guess)
	}

	jsGetSuggestedWords = js.FuncOf(func(_ js.Value, args []js.Value) any {
		if len(args) != 2 {
			return js.Global().Get("Error").New("Incorrect number of arguments to getSuggestedWords - must be 2")
		}
		guess := args[0].String()
		jsColourStringsArr := args[1]
		if !jsColourStringsArr.InstanceOf(js.Global().Get("Array")) {
			return js.Global().Get("Error").New("args[1] must be an array")
		}
		colourStringsLength := jsColourStringsArr.Length()
		if colourStringsLength != game.WordLength {
			return js.Global().Get("Error").New(fmt.Sprintf("colour pattern must be of length %d", game.WordLength))
		}
		var colourStringsSlice [game.WordLength]string
		for i := 0; i < len(colourStringsSlice); i++ {
			colourStringsSlice[i] = jsColourStringsArr.Index(i).String()
		}
		colourPattern, err := game.ColourStringsToColourPattern(colourStringsSlice)
		if err != nil {
			return js.Global().Get("Error").New(fmt.Sprintf("unable to parse colour strings: %s", err))
		}
		guessHelper.MakeGuess(guess, colourPattern)
		log.Printf("guesses: %s", guessHelper.Guesses)
		sortedGuessOutcomes := guessHelper.GetSortedGuessOutcomes(game.NormalMode)
		returnArr := js.Global().Get("Array").New()
		for _, guessOutcome := range sortedGuessOutcomes {
			returnArr.Call("push", guessOutcome.Guess)
		}
		return returnArr
	})

	jsUndoLastGuess = js.FuncOf(func(_ js.Value, args []js.Value) any {
		gh := guessHelper
		if gh == nil {
			return js.Global().Get("Error").New("guessHelper not yet initialised")
		}
		if err := guessHelper.RevertLastGuess(); err != nil && err != game.ErrNoGuesses {
			log.Printf("error: %s", err.Error())
			return js.Global().Get("Error").New(fmt.Sprintf("unable to revert last guess - err: %s", err))
		}
		return nil
	})

	jsResetGuessHelper = js.FuncOf(func(_ js.Value, args []js.Value) any {
		gh, err := game.NewGuessHelper(game.GuessHelperConfig{WordList: possibleAnswers, FreqMap: freqMap})
		if err != nil {
			js.Global().Get("Error").New(fmt.Sprintf("unable to create guessHelper - err: %s", err))
		}
		guessHelper = gh
		return nil
	})
	js.Global().Set("resetGuessHelper", jsResetGuessHelper)

	jsGuessHelper := js.Global().Get("Object").New()
	jsGuessHelper.Set("getSuggestedWords", jsGetSuggestedWords)
	jsGuessHelper.Set("undoLastGuess", jsUndoLastGuess)
	jsGuessHelper.Set("resetGuessHelper", jsResetGuessHelper)
	js.Global().Set("guessHelper", jsGuessHelper)

	js.Global().Set("allValidGuessesList", jsAllValidGuesses)
	js.Global().Set("optimalFirstGuesses", jsOptimalFirstGuesses)

	select {}
}
