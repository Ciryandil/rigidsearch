package indexing

import (
	"rigidsearch/stemming"
	"rigidsearch/stop_words"
	"rigidsearch/string_utils"
	"strings"
)

func ConstructTermFrequencyMap(text string) map[string]int {
	words := strings.Split(text, " ")
	freqCount := make(map[string]int)
	for _, word := range words {
		cleanedWord := string_utils.CleanWord(word)
		if _, ok := stop_words.STOP_WORDS[cleanedWord]; !ok {
			wordStem := stemming.PorterStemmer(cleanedWord)
			freqCount[wordStem] += 1
		}
	}
	return freqCount
}
