package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/hnasser-dev/wordle-solver/internal/game"
	"github.com/hnasser-dev/wordle-solver/internal/words"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type WordleHandler struct{}

type TodayGame struct {
	Answer  string   `json:"answer"`
	GameWon bool     `json:"gameWon"`
	Guesses []string `json:"guesses"`
}

func (h *WordleHandler) Today(c *echo.Context, wordList []string) error {
	answer := "hello"
	slog.Debug("answer", "answer", answer)
	mode := game.NormalMode
	initialGuesses := []string{"raise"}
	gameWon, guesses, err := playGame(answer, mode, initialGuesses, wordList)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	todayGame := TodayGame{
		Answer:  answer,
		GameWon: gameWon,
		Guesses: guesses,
	}
	return c.JSON(http.StatusOK, todayGame)
}

func playGame(answer string, mode game.GameMode, initialGuesses []string, wordList []string) (bool, []string, error) {
	gameWon := false
	guesses := []string{}
	game, err := game.NewGameSimulator(game.GameSimulatorConfig{Answer: answer, GameMode: mode, InitialGuesses: initialGuesses})
	if err != nil {
		slog.Error("unable to create new game", "err", err.Error())
		return gameWon, guesses, err
	}
	gameWon, guesses = game.PlayGameUntilEnd(true)
	slog.Info(
		"game complete",
		slog.String("answer", answer),
		slog.Int("numGuesses", len(guesses)),
		slog.Any("guesses", guesses),
		slog.Bool("gameWon", gameWon),
	)
	return gameWon, guesses, nil
}

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	wordList, err := words.GetWordList()
	if err != nil {
		slog.Error("unable to load word list", "error", err.Error())
		os.Exit(1)
	}

	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.Static("/", "web")

	wordleHandler := WordleHandler{}
	e.GET("/api/solvetoday", func(c *echo.Context) error { return wordleHandler.Today(c, wordList) })

	e.Start(":8080")
}
