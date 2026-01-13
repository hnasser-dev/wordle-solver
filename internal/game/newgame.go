package game

import "github.com/hnasser-dev/wordle-solver/internal/words"

type Game struct {
	GameWon                 bool
	Guesses                 []string
	InitialWordList         []string
	RemainingWordList       []string
	SortedRemainingOutcomes []guessOutcome
}

func NewGame() (*Game, error) {

	initialWordList, err := words.GetWordList()
	if err != nil {
		return nil, err
	}
	remainingWordList := make([]string, len(initialWordList))
	copy(remainingWordList, initialWordList)

	game := Game{
		GameWon:                 false,
		Guesses:                 []string{},
		InitialWordList:         initialWordList,
		RemainingWordList:       remainingWordList,
		SortedRemainingOutcomes: []guessOutcome{},
	}

	return &game, nil
}

func (g *Game) Guess(guess string) {
	// do a guess
}

func (g *Game) PlayGameUntilEnd(limitGuesses bool) {
	// play the game until complete
}
