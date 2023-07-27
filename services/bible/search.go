package bible

type SearchService interface {
}

type Storer interface {
}

type SearchRequest struct {
	Search Search `json:"search"`
}

type SearchResponse struct {
	SearchResults SearchResults `json:"search_results"`
	Error         string        `json:"error,omitempty"`
}
