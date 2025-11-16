package data_models

type Document struct {
	Name string
	Text string
	Id   string
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
