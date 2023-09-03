package bible

import (
	"context"

	"git.launchpad.net/~man4christ/+git/seed/auth"
	"git.launchpad.net/~man4christ/+git/seed/server"
	"github.com/kjvonly/service/services/bible/stores/elasticsearch"
	"go.uber.org/zap"
)

type BibleSearchService interface {
	Search(req BibleSearchRequest, gr server.GenericRequest) BibleSearchResponse
}

type Storer interface {
	Sql(ctx context.Context, sql string) (*elasticsearch.SqlResult, error)
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
	res, err := b.storer.Sql(gr.Ctx, req.Query)
	if err != nil {
		return BibleSearchResponse{
			Error: err.Error(),
		}
	}
	return BibleSearchResponse{
		SearchResults: SearchResults{SqlResults: *res},
	}
}

type BibleSearchRequest struct {
	Query string `json:"query"`
}

type BibleSearchResponse struct {
	SearchResults SearchResults `json:"search_results"`
	Error         string        `json:"error,omitempty"`
}

func (b BibleSearchServicer) Register(s *server.Server) {
	s.Register("BibleSearchService", "Search", server.RPCEndpoint{Roles: []string{}, Handler: b.SearchHandler})
}

// Create new BibleSearchServicer
func NewBibleSearchServicer(log *zap.SugaredLogger, storer Storer, a auth.Auth) BibleSearchRpcService {
	return BibleSearchServicer{
		log:    log,
		storer: storer,
		auth:   a,
	}
}
