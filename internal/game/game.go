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

const wordLength int = 5

const maxNumGuesses int = 6

const (
	Grey colour = iota
	Yellow
	Green
)

const (
	NormalMode gameMode = iota
	Dumb
)

var correctGuessColourPattern = colourPattern{Green, Green, Green, Green, Green}

type gameMode uint8

type colour uint8

type colourPattern [wordLength]colour

type guessDistribution map[colourPattern][]string

type guessOutcome struct {
	guess        string
	distribution guessDistribution
	entropyBits  float64
}

type Game struct {
	Answer                  string
	GameWon                 bool
	Guesses                 []string
	InitialWordList         []string
	RemainingWordList       []string
	SortedRemainingOutcomes []guessOutcome
	WordFrequencies         words.WordFrequencyMap
	GameMode                gameMode
}

func NewGame(answer string, mode gameMode, initialGuesses ...string) (*Game, error) {

	initialWordList, err := words.GetWordList()
	if err != nil {
		return nil, fmt.Errorf("unable to read word list - err: %w", err)
	}
	remainingWordList := make([]string, len(initialWordList))
	copy(remainingWordList, initialWordList)

	answerInWordList := slices.Contains(initialWordList, answer)
	if !answerInWordList {
		return nil, fmt.Errorf("provided answer %q is not in the word list", answer)
	}

	freqMap, err := words.GetWordFrequencyMap()
	if err != nil {
		return nil, fmt.Errorf("unable to read frequency map - err: %w", err)
	}

	game := Game{
		Answer:                  answer,
		GameWon:                 false,
		Guesses:                 []string{},
		InitialWordList:         initialWordList,
		RemainingWordList:       remainingWordList,
		SortedRemainingOutcomes: []guessOutcome{},
		WordFrequencies:         freqMap,
		GameMode:                mode,
	}

	for _, guess := range initialGuesses {
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
	case Dumb:
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
