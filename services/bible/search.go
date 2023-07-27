package bible

import (
	"github.com/gitamped/seed/auth"
	"github.com/gitamped/seed/server"
	"go.uber.org/zap"
)

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

type BibleSearchServicer struct {
	log    *zap.SugaredLogger
	storer Storer
	auth   auth.Auth
}

func (b BibleSearchServicer) Search(req BibleSearchRequest, gr server.GenericRequest) BibleSearchResponse {
	return BibleSearchResponse{}
}

type BibleSearchRequest struct {
	Search Search `json:"search"`
}

type BibleSearchResponse struct {
	SearchResults SearchResults `json:"search_results"`
	Error         string        `json:"error,omitempty"`
}
