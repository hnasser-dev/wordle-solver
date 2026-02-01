package words

import (
	_ "embed"
	"encoding/json"
	"strings"
)

//go:embed data/word_list.txt
var WordListData string

//go:embed data/word_frequency_map.json
var FreqMapData string

//go:embed data/optimal_first_guesses.txt
var OptimalFirstGuessesList string

type WordFrequencyMap map[string]float64

func GetWordList() []string {
	words := strings.Split(strings.TrimSpace(WordListData), "\n")
	return words
}

func GetWordFrequencyMap() (WordFrequencyMap, error) {
	var freqMap WordFrequencyMap
	err := json.Unmarshal([]byte(FreqMapData), &freqMap)
	if err != nil {
		return nil, err
	}
	return freqMap, nil
}

func GetOptimalFirstGuessesList() []string {
	words := strings.Split(strings.TrimSpace(OptimalFirstGuessesList), "\n")
	return words
}
