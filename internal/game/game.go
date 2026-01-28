package game

import (
	"fmt"
	"log/slog"
	"math"
	"reflect"
	"sort"

	"slices"

	"github.com/hnasser-dev/wordle-solver/internal/words"
)

const (
	wordLength    = 5
	maxNumGuesses = 6
)

const (
	Grey colour = iota
	Yellow
	Green
)

const (
	NormalMode GameMode = iota
	DumbMode
)

var correctGuessColourPattern = colourPattern{Green, Green, Green, Green, Green}

type GameMode uint8

func (m GameMode) Valid() bool {
	switch m {
	case NormalMode, DumbMode:
		return true
	default:
		return false
	}
}

type colour uint8
type colourPattern [wordLength]colour
type guessDistribution map[colourPattern][]string

type guessOutcome struct {
	guess        string
	distribution guessDistribution
	entropyBits  float64
}

type GameConfig struct {
	Answer   string
	GameMode GameMode

	InitialGuesses []string
	WordList       []string
	FreqMap        words.WordFrequencyMap
}

type Game struct {
	Answer                  string
	GameWon                 bool
	Guesses                 []string
	InitialWordList         []string
	RemainingWordList       []string
	SortedRemainingOutcomes []guessOutcome
	WordFrequencies         words.WordFrequencyMap
	GameMode                GameMode
}

func NewGame(config GameConfig) (*Game, error) {

	var err error

	if config.Answer == "" {
		return nil, fmt.Errorf("answer is required")
	}

	if !config.GameMode.Valid() {
		return nil, fmt.Errorf("GameMode not required or is invalid")
	}

	var initialWordList []string

	if config.WordList == nil {
		initialWordList, err = words.GetWordList()
		if err != nil {
			return nil, fmt.Errorf("unable to read word list - err: %w", err)
		}
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
		// copy
		for k, v := range config.FreqMap {
			freqMap[k] = v
		}
	}

	answerInWordList := slices.Contains(initialWordList, config.Answer)
	if !answerInWordList {
		return nil, fmt.Errorf("provided answer %q is not in the word list", config.Answer)
	}

	game := Game{
		Answer:                  config.Answer,
		GameWon:                 false,
		Guesses:                 []string{},
		InitialWordList:         initialWordList,
		RemainingWordList:       remainingWordList,
		SortedRemainingOutcomes: []guessOutcome{},
		WordFrequencies:         freqMap,
		GameMode:                config.GameMode,
	}

	for _, guess := range config.InitialGuesses {
		game.PerformGuess(guess)
	}

	return &game, nil
}

func (g *Game) PerformGuess(guess string) bool {
	slog.Debug("performing provided guess", slog.String("guess", guess))
	guessIsCorrect, remainingWordList := executeGuess(guess, g.Answer, g.RemainingWordList)
	g.Guesses = append(g.Guesses, guess)
	g.RemainingWordList = remainingWordList
	g.GameWon = guessIsCorrect || g.GameWon
	return g.GameWon
}

func (g *Game) PerformOptimalGuess() bool {
	sortedRemainingOutcomes := getSortedGuessOutcomes(g.RemainingWordList, g.WordFrequencies)
	var guess string
	switch g.GameMode {
	case DumbMode:
		guess = sortedRemainingOutcomes[len(sortedRemainingOutcomes)-1].guess
		slog.Debug("performing least optimal guess", "guess", guess)
	case NormalMode:
		guess = sortedRemainingOutcomes[0].guess
		slog.Debug("performing optimal guess", "guess", guess)
	default:
		panic(fmt.Sprintf("unknwon gameMode: %d", g.GameMode))
	}
	guessIsCorrect, remainingWordList := executeGuess(guess, g.Answer, g.RemainingWordList)
	g.Guesses = append(g.Guesses, guess)
	g.RemainingWordList = remainingWordList
	g.GameWon = guessIsCorrect || g.GameWon
	return g.GameWon
}

func (g *Game) PlayGameUntilEnd(limitGuesses bool) (bool, []string) {
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

	colourPatternSlice := [wordLength]colour{Grey, Grey, Grey, Grey, Grey}

	numCharsInAnswer := map[byte]int{}
	for i := range len(answer) {
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

	return colourPatternSlice
}

func computeGuessDistribution(guess string, wordList []string) guessDistribution {
	dist := guessDistribution{}
	for _, potentialAnswer := range wordList {
		if potentialAnswer == guess {
			continue
		}
		colourPattern := getColourPattern(guess, potentialAnswer)
		dist[colourPattern] = append(dist[colourPattern], potentialAnswer)
	}
	return dist
}

func getSortedGuessOutcomes(remainingWords []string, freqMap words.WordFrequencyMap) []guessOutcome {

	guessDistributions := map[string]guessDistribution{}
	for _, potentialGuess := range remainingWords {
		guessDistributions[potentialGuess] = computeGuessDistribution(potentialGuess, remainingWords)
	}

	guessOutcomes := make([]guessOutcome, 0, len(guessDistributions))
	for guess, dist := range guessDistributions {
		outcome := guessOutcome{
			guess:        guess,
			distribution: dist,
			entropyBits:  0.0,
		}
		for _, guesses := range dist {
			prob := float64(len(guesses)) / float64(len(remainingWords))
			outcome.entropyBits += float64(prob) * math.Log2(1.0/prob)
		}
		guessOutcomes = append(guessOutcomes, outcome)
	}

	sort.Slice(
		guessOutcomes,
		func(i, j int) bool {
			// if equal entropies, prioritise higher frequency
			if guessOutcomes[i].entropyBits == guessOutcomes[j].entropyBits {
				return freqMap[guessOutcomes[i].guess] > freqMap[guessOutcomes[j].guess]
			} else {
				return guessOutcomes[i].entropyBits > guessOutcomes[j].entropyBits
			}
		},
	)

	return guessOutcomes
}

// Internal helper function that applies the guess (mutates Game)
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
