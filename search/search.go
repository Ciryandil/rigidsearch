package search

import (
	"fmt"
	"math"
	"rigidsearch/data_models"
	"rigidsearch/heap"
	"rigidsearch/indexing"
	"rigidsearch/stemming"
	"rigidsearch/stop_words"
	"strings"
)

func Search(query data_models.Query) ([]data_models.SearchResult, error) {
	queryTerms := strings.Split(query.Query, " ")
	finalQueryTerms := make([]string, 0)
	for _, term := range queryTerms {
		if _, ok := stop_words.STOP_WORDS[term]; !ok {
			finalQueryTerms = append(finalQueryTerms, stemming.PorterStemmer(term))
		}
	}
	if len(finalQueryTerms) == 0 {
		return nil, fmt.Errorf("all terms in query were stop words! Use non stop words")
	}
	if query.NumResults == 0 {
		query.NumResults = 5
	}
	idfDenom := float64(1 + indexing.GlobalSearchIndex.TotalDocs)
	docsMap := make(map[string]float64)
	indexing.GlobalSearchIndex.Lock.RLock()
	defer indexing.GlobalSearchIndex.Lock.RUnlock()
	for _, term := range finalQueryTerms {
		freqData, ok := indexing.GlobalSearchIndex.WordFrequencyMap[term]
		if !ok {
			continue
		}
		var docFreq int
		docFreqData := indexing.GlobalSearchIndex.WordToDocMap[term]
		if docFreqData != nil {
			docFreq = docFreqData.TotalDocs
		}
		inverseDocFrequency := float64(docFreq)/idfDenom + 1
		inverseDocFrequency = math.Log(inverseDocFrequency)
		tfDenom := float64(freqData.TotalFrequency)
		for docId, freq := range freqData.FrequencyMap {
			termFreq := float64(freq) / tfDenom
			score := termFreq * inverseDocFrequency
			docsMap[docId] += score
		}
	}

	resultArr := make([]data_models.IntermediateResult, 0)
	for docId, score := range docsMap {
		resultArr = append(resultArr, data_models.IntermediateResult{DocId: docId, Score: score})
	}
	heap.Heapify(resultArr, data_models.IntermediateResultComparator)
	topResults := make([]data_models.SearchResult, 0)
	for itr := 0; itr < query.NumResults; itr += 1 {
		resPtr := heap.Pop(resultArr, data_models.IntermediateResultComparator)
		if resPtr == nil {
			break
		}
		docData := indexing.GlobalSearchIndex.DocMetadataMap[resPtr.DocId]
		topResults = append(topResults, data_models.SearchResult{DocId: docData.Id, Name: docData.Name, Score: resPtr.Score})
	}
	return topResults, nil
}
