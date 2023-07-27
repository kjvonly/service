package bible

import "github.com/gitamped/seed/server"

type BibleSearchService interface {
	Search(BibleSearchRequest) BibleSearchResponse
}

type Storer interface {
}

// Required to register endpoints with the Server
type BibleSearchRpcService interface {
	BibleSearchService
	// Registers RPCService with Server
	Register(s *server.Server)
}

type BibleSearchRequest struct {
	Search Search `json:"search"`
}

type BibleSearchResponse struct {
	SearchResults SearchResults `json:"search_results"`
	Error         string        `json:"error,omitempty"`
}
