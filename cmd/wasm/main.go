package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func getDateNowStr() string {
	d := js.Global().Get("Date").New()
	year := d.Call("getFullYear").Int()
	month := d.Call("getMonth").Int() + 1 // JS months are 0-11
	day := d.Call("getDate").Int()
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

func getAnswerFromApi() (string, error) {
	dateNowStr := getDateNowStr()
	resp, err := http.Get(fmt.Sprintf("https://www.nytimes.com/svc/wordle/v2/%s.json", dateNowStr))
	if err != nil {
		return "", fmt.Errorf("unable to get answer from nytimes api - err: %w", err)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	data := nyTimesApiResponse{}
	err = decoder.Decode(&data)
	if err != nil {
		return "", fmt.Errorf("unable to decode response - err: %w", err)
	}
	return data.Solution, nil
}

func solveWordle(this js.Value, args []js.Value) interface{} {
	gameModeStr := args[0]
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
	if len(args) > 1 {
		for i, arg := range args[1:] {
			initialGuesses[i] = arg.String()
		}
	}
	answer, err := getAnswerFromApi()
	if err != nil {
		return js.ValueOf(fmt.Sprintf("error: unable to retrieve answer from API - err: %s", err.Error()))
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
	jsObj.Set("guesses", guesses)
	return jsObj
}

func main() {
	wordList, _ = words.GetWordList()
	freqMap, _ = words.GetWordFrequencyMap()
	js.Global().Set("solveWordle", js.FuncOf(solveWordle))
}
