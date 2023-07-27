package bible

import "github.com/kjvonly/service/services/bible/stores/elasticsearch"

type Search struct {
	Query string `json:"query"`
}
type SearchResults struct {
	SqlResults elasticsearch.SqlResult
}
