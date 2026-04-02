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

type SynchronousIndex struct {
	Index SearchIndex
	Lock  sync.RWMutex
}

var GlobalSearchIndex SynchronousIndex

type Posting struct {
	DocId int32
	Tf    int32 // term frequency - frequency of word in doc of DocId
}

type TermInfo struct {
	Postings       []Posting
	DocFrequency   int32 // number of docs in which word appears
	TotalFrequency int32 // total frequency of word across docs
}

type SearchIndex struct {
	TermIndex      map[string]int32                       `json:"term_index"` // map from term to integer index of it in terms array
	Terms          []*TermInfo                            `json:"terms"`      // array of terminfo
	DocMetadataMap map[int32]data_models.DocumentMetadata `json:"doc_metadata_map"`
	DeletedDocs    map[int32]struct{}                     `json:"deleted_docs"` // ids of docs that have been deleted - lazy deletion
}

func LoadIndex() error {
	GlobalSearchIndex.Lock.Lock()
	defer GlobalSearchIndex.Lock.Unlock()
	var index SearchIndex
	index.DocMetadataMap = make(map[int32]data_models.DocumentMetadata)
	index.Terms = nil
	index.TermIndex = make(map[string]int32)
	GlobalSearchIndex.Index = index
	file, err := os.ReadFile(constants.INDEX_FILE)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	var storedIndex SearchIndex
	err = json.Unmarshal(file, &storedIndex)
	if err != nil {
		return err
	}
	if len(storedIndex.Terms) != len(storedIndex.TermIndex) {
		return fmt.Errorf("error loading index file, word and their frequencies don't match in lengths")
	}

	GlobalSearchIndex.Index = storedIndex
	return nil
}

func StoreIndex() error {
	GlobalSearchIndex.Lock.Lock()
	defer GlobalSearchIndex.Lock.Unlock()

	bytes, err := json.Marshal(GlobalSearchIndex.Index)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(constants.INDEX_FILE, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(bytes)
	return err
}
