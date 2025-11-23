package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"rigidsearch/constants"
	"rigidsearch/data_models"
	"rigidsearch/indexing"
	"rigidsearch/search"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func returnError(w http.ResponseWriter, err error, statusCode int) {
	http.Error(w, fmt.Sprintf("%v", err), statusCode)
}

func NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Operational"))
	})
	r.Get("/search", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		query := q.Get("query")
		numResultsStr := q.Get("num_results")
		method := q.Get("method")
		if method == "" {
			method = "tf_idf"
		}
		numResults, err := strconv.Atoi(numResultsStr)
		if err != nil {
			numResults = 0
		}
		queryStruct := data_models.Query{
			Query:      query,
			NumResults: numResults,
		}
		var results []data_models.SearchResult
		if method == "bm_25" {
			results, err = search.Bm25Search(queryStruct)
		} else {
			results, err = search.TfIdfSearch(queryStruct)
		}
		if err != nil {
			returnError(w, err, http.StatusInternalServerError)
			return
		}
		jsonResults, err := json.Marshal(results)
		if err != nil {
			returnError(w, err, http.StatusInternalServerError)
			return
		}
		w.Write(jsonResults)
	})
	r.Post("/index", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			returnError(w, fmt.Errorf("invalid content-type, expected application/json"), http.StatusUnsupportedMediaType)
			return
		}
		var indexRequest data_models.Document
		err := json.NewDecoder(r.Body).Decode(&indexRequest)
		if err != nil {
			returnError(w, fmt.Errorf("failed to decode request: %w", err), http.StatusBadRequest)
			return
		}
		docId, err := indexing.IndexDocument(indexRequest)
		if err != nil {
			returnError(w, err, http.StatusInternalServerError)
			return
		}
		result := map[string]interface{}{
			"document_id": docId,
		}
		resp, err := json.Marshal(result)
		if err != nil {
			returnError(w, err, http.StatusInternalServerError)
			return
		}
		w.Write(resp)
	})
	r.Delete("/documents/{documentId}", func(w http.ResponseWriter, r *http.Request) {
		docId := chi.URLParam(r, "documentId")
		err := indexing.DeleteDocument(docId)
		if err != nil {
			returnError(w, err, http.StatusInternalServerError)
			return
		}
		w.Write([]byte("success"))
	})
	r.Get("/documents/{documentId}", func(w http.ResponseWriter, r *http.Request) {
		docId := chi.URLParam(r, "documentId")
		bytes, err := os.ReadFile(fmt.Sprintf("%s/%s", constants.STORAGE_LOC, docId))
		if err != nil {
			returnError(w, err, http.StatusInternalServerError)
		}
		docMetadata := indexing.GlobalSearchIndex.DocMetadataMap[docId]
		result := map[string]interface{}{
			"name": docMetadata.Name,
			"text": string(bytes),
		}
		resp, err := json.Marshal(result)
		if err != nil {
			returnError(w, err, http.StatusInternalServerError)
			return
		}
		w.Write(resp)
	})
	return r
}
