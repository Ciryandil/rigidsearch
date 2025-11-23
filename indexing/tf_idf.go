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
	GlobalSearchIndex.Lock.Lock()
	defer GlobalSearchIndex.Lock.Unlock()
	termFrequencyMap := ConstructTermFrequencyMap(document.Text)
	docId := uuid.NewString()
	f, err := os.OpenFile(fmt.Sprintf("%s/%s", constants.STORAGE_LOC, docId), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = f.Write([]byte(document.Text))
	if err != nil {
		return "", err
	}
	docLength := 0
	for word, freq := range termFrequencyMap {
		docLength += freq
		GlobalSearchIndex.DocToWordMap[docId] = append(GlobalSearchIndex.DocToWordMap[docId], word)
		if _, ok := GlobalSearchIndex.WordFrequencyMap[word]; !ok {
			GlobalSearchIndex.WordFrequencyMap[word] = &WordFrequencyData{
				FrequencyMap:   make(map[string]int),
				TotalFrequency: 0,
			}
		}
		GlobalSearchIndex.WordFrequencyMap[word].FrequencyMap[docId] = freq
		GlobalSearchIndex.WordFrequencyMap[word].TotalFrequency += freq
		freqData := GlobalSearchIndex.WordToDocMap[word]
		if freqData == nil {
			freqData = &DocFrequencyData{
				DocSet: make(map[string]struct{}),
			}
		}
		freqData.DocSet[docId] = struct{}{}
		GlobalSearchIndex.WordToDocMap[word] = freqData
	}
	document.Id = docId
	document.Length = docLength
	GlobalSearchIndex.DocMetadataMap[docId] = document
	return docId, nil
}

func DeleteDocument(documentId string) error {
	GlobalSearchIndex.Lock.Lock()
	defer GlobalSearchIndex.Lock.Unlock()
	words := GlobalSearchIndex.DocToWordMap[documentId]
	fmt.Println("Doc id: ", documentId)
	for _, word := range words {
		freqData := GlobalSearchIndex.WordFrequencyMap[word]
		if freqData != nil {
			count := freqData.FrequencyMap[documentId]
			delete(freqData.FrequencyMap, documentId)
			freqData.TotalFrequency -= count
			if len(freqData.FrequencyMap) == 0 {
				delete(GlobalSearchIndex.WordFrequencyMap, word)
			}
			fmt.Println("Word: ", word, " after deletion freq data: ", freqData.FrequencyMap)
		}
		docFreqData := GlobalSearchIndex.WordToDocMap[word]
		if docFreqData != nil {
			delete(docFreqData.DocSet, documentId)
			delete(GlobalSearchIndex.WordToDocMap, word)
		}
	}
	delete(GlobalSearchIndex.DocToWordMap, documentId)
	delete(GlobalSearchIndex.DocMetadataMap, documentId)

	err := os.Remove(fmt.Sprintf("%s/%s", constants.STORAGE_LOC, documentId))
	return err

}
