/*
Struct for simulating optimal (or suboptimal) wordle games where the answer is known ahead of time
*/

package game

import (
	"fmt"
	"log/slog"
	"maps"
	"reflect"

	"slices"

	"github.com/hnasser-dev/wordle-solver/internal/words"
)

const (
	WordLength    = 5
	maxNumGuesses = 6
)

var correctGuessColourPattern = colourPattern{Green, Green, Green, Green, Green}

type GameSimulatorConfig struct {
	Answer   string
	GameMode GameMode

	InitialGuesses []string
	WordList       []string
	FreqMap        words.WordFrequencyMap
}

type GameSimulator struct {
	Answer                  string
	GameWon                 bool
	Guesses                 []string
	InitialWordList         []string
	RemainingWordList       []string
	SortedRemainingOutcomes []guessOutcome
	FreqMap                 words.WordFrequencyMap
	GameMode                GameMode
}

func NewGameSimulator(config GameSimulatorConfig) (*GameSimulator, error) {
	var err error
	if config.Answer == "" {
		return nil, fmt.Errorf("answer is required")
	}
	if !config.GameMode.Valid() {
		return nil, fmt.Errorf("GameMode not provided or is invalid")
	}
	var initialWordList []string
	if config.WordList == nil {
		initialWordList = words.GetWordList()
	} else {
		initialWordList = append([]string{}, config.WordList...) // copy
	}
	remainingWordList := append([]string{}, initialWordList...) // copy
	freqMap := words.WordFrequencyMap{}
	if config.FreqMap == nil {
		freqMap, err = words.GetWordFrequencyMap()
		if err != nil {
			return nil, fmt.Errorf("unable to read frequency map - err: %w", err)
		}
	} else {
		maps.Copy(freqMap, config.FreqMap)
	}
	answerInWordList := slices.Contains(initialWordList, config.Answer)
	if !answerInWordList {
		return nil, fmt.Errorf("provided answer %q is not in the word list", config.Answer)
	}
	game := GameSimulator{
		Answer:                  config.Answer,
		GameWon:                 false,
		Guesses:                 []string{},
		InitialWordList:         initialWordList,
		RemainingWordList:       remainingWordList,
		SortedRemainingOutcomes: []guessOutcome{},
		FreqMap:                 freqMap,
		GameMode:                config.GameMode,
	}
	for _, guess := range config.InitialGuesses {
		game.PerformGuess(guess)
	}
	return &game, nil
}

func (g *GameSimulator) PerformGuess(guess string) bool {
	slog.Debug("performing provided guess", slog.String("guess", guess))
	guessIsCorrect, remainingWordList := executeGuess(guess, g.Answer, g.RemainingWordList)
	g.Guesses = append(g.Guesses, guess)
	g.RemainingWordList = remainingWordList
	g.GameWon = guessIsCorrect || g.GameWon
	return g.GameWon
}

func (g *GameSimulator) PerformOptimalGuess() bool {
	sortedRemainingOutcomes := getSortedGuessOutcomes(g.RemainingWordList, g.FreqMap)
	var guess string
	switch g.GameMode {
	case DumbMode:
		guess = sortedRemainingOutcomes[len(sortedRemainingOutcomes)-1].Guess
		slog.Debug("performing least optimal guess", "guess", guess)
	case NormalMode:
		guess = sortedRemainingOutcomes[0].Guess
		slog.Debug("performing optimal guess", "guess", guess)
	default:
		panic(fmt.Sprintf("unknown gameMode: %d", g.GameMode))
	}
	guessIsCorrect, remainingWordList := executeGuess(guess, g.Answer, g.RemainingWordList)
	g.Guesses = append(g.Guesses, guess)
	g.RemainingWordList = remainingWordList
	g.GameWon = guessIsCorrect || g.GameWon
	return g.GameWon
}

func (g *GameSimulator) PlayGameUntilEnd(limitGuesses bool) (bool, []string) {
	// play the game until complete
	for !g.GameWon {
		if limitGuesses && len(g.Guesses) == maxNumGuesses {
			break
		}
		g.PerformOptimalGuess()
	}
	return g.GameWon, g.Guesses
}

func getColourPattern(guess string, answer string) colourPattern {
	colourPatternSlice := [WordLength]colour{Grey, Grey, Grey, Grey, Grey}
	numCharsInAnswer := map[byte]int{}
	for i := range len(answer) {
		numCharsInAnswer[answer[i]]++
	}
	numCharsGuessed := map[byte]int{}
	// assign greens
	for idx := range WordLength {
		if guess[idx] == answer[idx] {
			colourPatternSlice[idx] = Green
			numCharsGuessed[guess[idx]]++
		}
	}
	// assign yellows
	for idx := range WordLength {
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

func executeGuess(guess string, answer string, remainingWords []string) (bool, []string) {
	guessIsCorrect := false
	guessDistribution := computeGuessDistribution(guess, remainingWords)
	colourPattern := getColourPattern(guess, answer)
	if reflect.DeepEqual(colourPattern, correctGuessColourPattern) {
		guessIsCorrect = true
	}
	resultngWordList := guessDistribution[colourPattern]
	return guessIsCorrect, resultngWordList
}
