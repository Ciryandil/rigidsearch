package data_models

type Document struct {
	Name   string `json:"name"`
	Text   string `json:"text"`
	Id     string `json:"id"`
	Length int    `json:"length"`
}

type Query struct {
	Query      string
	NumResults int
}

type SearchResult struct {
	DocId string  `json:"id"`
	Name  string  `json:"name"`
	Score float64 `json:"score"`
}

type IntermediateResult struct {
	DocId string
	Score float64
}
