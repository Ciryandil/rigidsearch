package search

import (
	"fmt"
	"math"
	"rigidsearch/data_models"
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
		return nil, fmt.Errorf("All terms in query were stop words! Use non stop words")
	}
	if query.NumResults == 0 {
		query.NumResults = 5
	}
	idfDenom := float64(1 + indexing.GlobalSearchIndex.TotalDocs)
	docsMap := make(map[string]float64)
	for _, term := range finalQueryTerms {
		freqData, ok := indexing.GlobalSearchIndex.WordFrequencyMap[term]
		if !ok {
			continue
		}
		inverseDocFrequency := float64(len(indexing.GlobalSearchIndex.WordToDocMap[term]))/idfDenom + 1
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

}
