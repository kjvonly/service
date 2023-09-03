package elasticsearch

type SqlResult struct {
	Columns []Column `json:"columns"`
	Rows    [][]any  `json:"rows"`
}
type Column struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
