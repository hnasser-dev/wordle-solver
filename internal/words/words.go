package words

import (
	"fmt"
	"os"
	"strings"
)

const wordListFilePath string = "data/word_list.txt"

func GetWordList() ([]string, error) {
	wordList := []string{}
	dat, err := os.ReadFile(wordListFilePath)
	if err != nil {
		return wordList, fmt.Errorf("unable to read word list - err: %s", err)
	}
	wordList = strings.Split(string(dat), "\n")
	return wordList, nil
}
