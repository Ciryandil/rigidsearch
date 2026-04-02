package data_models

type Document struct {
	Name   string `json:"name"`
	Text   string `json:"text"`
	Id     int32  `json:"id"`
	Length int32  `json:"length"`
}

type DocumentMetadata struct {
	Name   string `json:"name"`
	Id     int32  `json:"id"`
	Length int32  `json:"length"`
}

type Query struct {
	Query      string
	NumResults int
}

type SearchResult struct {
	DocId int32   `json:"id"`
	Name  string  `json:"name"`
	Score float64 `json:"score"`
}

type IntermediateResult struct {
	DocId int32
	Score float64
}
