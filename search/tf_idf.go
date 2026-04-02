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

func TfIdfSearch(query data_models.Query) ([]data_models.SearchResult, error) {
	queryTerms := strings.Fields(query.Query)
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
	idfNom := float64(len(indexing.GlobalSearchIndex.Index.DocMetadataMap))
	docsMap := make(map[int32]float64)
	indexing.GlobalSearchIndex.Lock.RLock()
	defer indexing.GlobalSearchIndex.Lock.RUnlock()
	for _, term := range finalQueryTerms {
		termIndex, ok := indexing.GlobalSearchIndex.Index.TermIndex[term]
		if !ok {
			continue
		}
		termInfo := indexing.GlobalSearchIndex.Index.Terms[termIndex]
		if termInfo == nil {
			continue
		}

		docFreq := termInfo.DocFrequency
		inverseDocFrequency := idfNom/(float64(docFreq)+1) + 1
		fmt.Println("Term: ", term, " idf before log: ", inverseDocFrequency)
		inverseDocFrequency = math.Log(inverseDocFrequency)
		for _, posting := range termInfo.Postings {
			_, deleted := indexing.GlobalSearchIndex.Index.DeletedDocs[posting.DocId]
			if deleted {
				continue
			}
			docMetadata := indexing.GlobalSearchIndex.Index.DocMetadataMap[posting.DocId]
			docLength := docMetadata.Length
			if docLength == 0 {
				continue
			}
			termFreq := float64(posting.Tf) / float64(docLength)
			score := termFreq * inverseDocFrequency
			fmt.Println("Term: ", term, " score: ", score)

			docsMap[posting.DocId] += score
		}
	}

	resultArr := make([]data_models.IntermediateResult, 0)
	for docId, score := range docsMap {
		fmt.Println("Doc id: ", docId, " score: ", score)
		resultArr = append(resultArr, data_models.IntermediateResult{DocId: docId, Score: score})
	}
	heap.Heapify(resultArr, data_models.IntermediateResultComparator)
	topResults := make([]data_models.SearchResult, 0)
	for itr := 0; itr < query.NumResults; itr += 1 {
		fmt.Printf("Current heap: %v\n", resultArr)
		var resPtr *data_models.IntermediateResult
		resPtr, resultArr = heap.Pop(resultArr, data_models.IntermediateResultComparator)
		if resPtr == nil {
			break
		}
		docData := indexing.GlobalSearchIndex.Index.DocMetadataMap[resPtr.DocId]
		topResults = append(topResults, data_models.SearchResult{DocId: docData.Id, Name: docData.Name, Score: resPtr.Score})
	}
	return topResults, nil
}
