package bible

type Search struct {
	Query string `json:"query"`
}
type SearchResults struct {
	Columns []Column `json:"columns"`
	Rows    [][]any  `json:"rows"`
}

type Column struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
