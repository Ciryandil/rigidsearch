package indexing

import (
	"fmt"
	"math"
	"os"
	"rigidsearch/constants"
	"rigidsearch/data_models"
	"rigidsearch/stemming"
	"rigidsearch/stop_words"
	"rigidsearch/string_utils"
	"strings"
)

func ConstructTermFrequencyMap(text string) map[string]int {
	words := strings.Fields(text)
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

func IndexDocument(document data_models.Document) (int32, error) {
	GlobalSearchIndex.Lock.Lock()
	defer GlobalSearchIndex.Lock.Unlock()
	termFrequencyMap := ConstructTermFrequencyMap(document.Text)
	docId := int32(len(GlobalSearchIndex.Index.DocMetadataMap))
	f, err := os.OpenFile(fmt.Sprintf("%s/%d", constants.STORAGE_LOC, docId), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return -1, err
	}
	defer f.Close()
	_, err = f.Write([]byte(document.Text))
	if err != nil {
		return -1, err
	}
	for word, freq := range termFrequencyMap {
		document.Length += int32(freq)
		index, ok := GlobalSearchIndex.Index.TermIndex[word]
		if !ok {
			totalTerms := len(GlobalSearchIndex.Index.TermIndex)
			GlobalSearchIndex.Index.TermIndex[word] = int32(totalTerms)
			termInfo := TermInfo{
				Postings: []Posting{
					Posting{
						DocId: docId,
						Tf:    int32(freq),
					},
				},
				DocFrequency:   1,
				TotalFrequency: int32(freq),
			}
			GlobalSearchIndex.Index.Terms = append(GlobalSearchIndex.Index.Terms, &termInfo)
		} else {
			termInfo := GlobalSearchIndex.Index.Terms[index]
			termInfo.Postings = append(termInfo.Postings, Posting{
				DocId: docId,
				Tf:    int32(freq),
			})
			termInfo.DocFrequency += 1
			termInfo.TotalFrequency += int32(freq)
		}
	}
	document.Id = docId
	GlobalSearchIndex.Index.DocMetadataMap[docId] = data_models.DocumentMetadata{
		Id:     document.Id,
		Name:   document.Name,
		Length: document.Length,
	}
	return docId, nil
}

func DeleteDocument(documentId int) error {
	if documentId > math.MaxInt32 || documentId < math.MinInt32 {
		return fmt.Errorf("invalid document id")
	}
	documentIdI32 := int32(documentId)
	GlobalSearchIndex.Lock.Lock()
	defer GlobalSearchIndex.Lock.Unlock()
	delete(GlobalSearchIndex.Index.DocMetadataMap, documentIdI32)
	GlobalSearchIndex.Index.DeletedDocs[documentIdI32] = struct{}{}

	err := os.Remove(fmt.Sprintf("%s/%d", constants.STORAGE_LOC, documentId))
	return err

}
