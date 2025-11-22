package indexing

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"rigidsearch/constants"
	"rigidsearch/data_models"
	"sync"
)

var GlobalSearchIndex SearchIndex

type WordFrequencyData struct {
	FrequencyMap   map[string]int
	TotalFrequency int
}

type DocFrequencyData struct {
	DocSet    map[string]struct{}
	TotalDocs int
}

type SearchIndex struct {
	Lock             sync.RWMutex
	WordFrequencyMap map[string]*WordFrequencyData
	WordToDocMap     map[string]*DocFrequencyData
	DocToWordMap     map[string][]string
	DocMetadataMap   map[string]data_models.Document
	TotalDocs        int
}

type IndexData struct {
	Words           []string                        `json:"words"`
	WordFrequencies []map[string]int                `json:"word_frequencies"`
	DocLists        [][]string                      `json:"doc_lists"`
	DocMetadataMap  map[string]data_models.Document `json:"doc_metadata_map"`
}

func LoadIndex() error {
	GlobalSearchIndex.Lock.Lock()
	defer GlobalSearchIndex.Lock.Unlock()
	GlobalSearchIndex.DocMetadataMap = make(map[string]data_models.Document)
	GlobalSearchIndex.DocToWordMap = make(map[string][]string)
	GlobalSearchIndex.WordFrequencyMap = make(map[string]*WordFrequencyData)
	GlobalSearchIndex.WordToDocMap = make(map[string]*DocFrequencyData)
	file, err := os.ReadFile(constants.INDEX_FILE)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	if errors.Is(err, os.ErrNotExist) {
		GlobalSearchIndex.DocMetadataMap = make(map[string]data_models.Document)
		return nil
	}
	var indexData IndexData
	err = json.Unmarshal(file, &indexData)
	if err != nil {
		return err
	}
	if len(indexData.Words) != len(indexData.WordFrequencies) || len(indexData.Words) != len(indexData.DocLists) {
		return fmt.Errorf("error loading index file, word and their frequencies don't match in lengths")
	}
	for itr := range indexData.Words {
		frequencyData := WordFrequencyData{
			FrequencyMap:   make(map[string]int),
			TotalFrequency: 0,
		}
		docFrequencyData := DocFrequencyData{
			DocSet:    make(map[string]struct{}),
			TotalDocs: 0,
		}
		for doc, freq := range indexData.WordFrequencies[itr] {
			frequencyData.FrequencyMap[doc] = freq
			frequencyData.TotalFrequency += freq
			GlobalSearchIndex.DocToWordMap[doc] = append(GlobalSearchIndex.DocToWordMap[doc], indexData.Words[itr])
		}
		for _, doc := range indexData.DocLists[itr] {
			docFrequencyData.DocSet[doc] = struct{}{}
			docFrequencyData.TotalDocs += 1
		}
		GlobalSearchIndex.WordFrequencyMap[indexData.Words[itr]] = &frequencyData
		GlobalSearchIndex.WordToDocMap[indexData.Words[itr]] = &docFrequencyData
	}
	GlobalSearchIndex.DocMetadataMap = indexData.DocMetadataMap
	GlobalSearchIndex.TotalDocs = len(indexData.DocMetadataMap)
	return nil
}

func StoreIndex() error {
	GlobalSearchIndex.Lock.Lock()
	defer GlobalSearchIndex.Lock.Unlock()
	var indexData IndexData
	var wordList []string
	for word := range GlobalSearchIndex.WordFrequencyMap {
		wordList = append(wordList, word)
	}
	for _, word := range wordList {
		indexData.Words = append(indexData.Words, word)
		indexData.WordFrequencies = append(indexData.WordFrequencies, GlobalSearchIndex.WordFrequencyMap[word].FrequencyMap)
		docList := make([]string, 0)
		freqData := GlobalSearchIndex.WordToDocMap[word]
		if freqData != nil {
			for docId := range freqData.DocSet {
				docList = append(docList, docId)
			}
		}
		indexData.DocLists = append(indexData.DocLists, docList)
	}
	indexData.DocMetadataMap = GlobalSearchIndex.DocMetadataMap
	bytes, err := json.Marshal(indexData)
	if err != nil {
		return err
	}
	err = os.WriteFile(constants.INDEX_FILE, bytes, 0644)
	return err
}
