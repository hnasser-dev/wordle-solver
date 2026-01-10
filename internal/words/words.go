package words

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const wordListFilePath string = "data/word_list.txt"
const wordFreqMapFilePath string = "data/word_frequency_map.json"

type WordFrequencyMap map[string]float64

func GetWordList() ([]string, error) {
	wordList := []string{}
	dat, err := os.ReadFile(wordListFilePath)
	if err != nil {
		return wordList, fmt.Errorf("unable to read word list - err: %s", err)
	}
	wordList = strings.Split(string(dat), "\n")
	return wordList, nil
}

func GetWordFrequencyMap() (WordFrequencyMap, error) {
	data, err := os.ReadFile(wordFreqMapFilePath)
	if err != nil {
		return nil, err
	}
	var freqMap WordFrequencyMap
	if err := json.Unmarshal(data, &freqMap); err != nil {
		return nil, err
	}
	return freqMap, nil
}
