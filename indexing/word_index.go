package indexing

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"rigidsearch/constants"
	"rigidsearch/data_models"
)

var GlobalSearchIndex SearchIndex

type SearchIndex struct {
	WordFrequencyMap map[string]map[string]int
	WordToDocMap     map[string][]string
	DocMetadataMap   map[string]data_models.Document
}

type IndexData struct {
	Words           []string                        `json:"words"`
	WordFrequencies []map[string]int                `json:"word_frequencies"`
	DocLists        [][]string                      `json:"doc_lists"`
	DocMetadataMap  map[string]data_models.Document `json:"doc_metadata_map"`
}

func LoadIndex() error {
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
		GlobalSearchIndex.WordFrequencyMap[indexData.Words[itr]] = indexData.WordFrequencies[itr]
		GlobalSearchIndex.WordToDocMap[indexData.Words[itr]] = indexData.DocLists[itr]
	}
	GlobalSearchIndex.DocMetadataMap = indexData.DocMetadataMap
	return nil
}

func StoreIndex() error {
	var indexData IndexData
	var wordList []string
	for word := range GlobalSearchIndex.WordFrequencyMap {
		wordList = append(wordList, word)
	}
	for _, word := range wordList {
		indexData.Words = append(indexData.Words, word)
		indexData.WordFrequencies = append(indexData.WordFrequencies, GlobalSearchIndex.WordFrequencyMap[word])
		indexData.DocLists = append(indexData.DocLists, GlobalSearchIndex.WordToDocMap[word])
	}
	indexData.DocMetadataMap = GlobalSearchIndex.DocMetadataMap
	bytes, err := json.Marshal(indexData)
	if err != nil {
		return err
	}
	err = os.WriteFile(constants.INDEX_FILE, bytes, 0644)
	return err
}
