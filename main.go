package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const wordListFilePath string = "data/word_list.txt"

func getFullWordList() ([]string, error) {
	wordList := []string{}
	dat, err := os.ReadFile(wordListFilePath)
	if err != nil {
		return []string{}, fmt.Errorf("unable to read word list - err: %s", err)
	}
	wordList = strings.Split(string(dat), "\n")
	return wordList, nil
}

func playGame(wordList []string, answer string) error {
	return nil
}

func main() {
	wordList, err := getFullWordList()
	if err != nil {
		log.Fatal(err)
	}

	// for testing
	answer := "fruit"

	// TODO
	playGame(wordList, answer)
}
