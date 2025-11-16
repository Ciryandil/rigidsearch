package indexing

import (
	"fmt"
	"os"
	"rigidsearch/constants"
	"rigidsearch/data_models"
	"rigidsearch/stemming"
	"rigidsearch/stop_words"
	"rigidsearch/string_utils"
	"strings"

	"github.com/google/uuid"
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

func IndexDocument(document data_models.Document) (string, error) {
	termFrequencyMap := ConstructTermFrequencyMap(document.Text)
	docId := uuid.NewString()
	err := os.WriteFile(fmt.Sprintf("%s/%s", constants.STORAGE_LOC, docId), []byte(document.Text), 0644)
	if err != nil {
		return "", err
	}
	for word, freq := range termFrequencyMap {
		if _, ok := GlobalSearchIndex.WordFrequencyMap[word]; !ok {
			GlobalSearchIndex.WordFrequencyMap[word] = &WordFrequencyData{
				FrequencyMap:   make(map[string]int),
				TotalFrequency: 0,
			}
		}
		GlobalSearchIndex.WordFrequencyMap[word].FrequencyMap[docId] = freq
		GlobalSearchIndex.WordFrequencyMap[word].TotalFrequency += freq
		GlobalSearchIndex.WordToDocMap[word] = append(GlobalSearchIndex.WordToDocMap[word], docId)
	}
	GlobalSearchIndex.DocMetadataMap[docId] = document
	GlobalSearchIndex.TotalDocs += 1
	return docId, nil
}
