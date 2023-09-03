package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type Store struct {
	log              *zap.SugaredLogger
	elasticSearchUrl string
}

func (s Store) Sql(ctx context.Context, sql string) (*SqlResult, error) {
	requestURL := fmt.Sprintf("%s/_sql?format=json", s.elasticSearchUrl)
	b, err := json.Marshal(struct {
		Query string `json:"query"`
	}{Query: sql})

	if err != nil {

	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(b))

	if err != nil {
		errMsg := fmt.Sprintf("client: could not create request: %s\n", err)
		s.log.Infof(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		errMsg := fmt.Sprintf("client: error making http request: %s\n", err)
		s.log.Infof(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	if res.StatusCode != 200 {
		errMsg := fmt.Sprintf("client: %d http response\n", res.StatusCode)
		s.log.Infof(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		errMsg := fmt.Sprintf("client: could not read response body: %s\n", err)
		s.log.Infof(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	var sqlResult SqlResult
	err = json.Unmarshal(resBody, &sqlResult)
	if err != nil {
		errMsg := fmt.Sprintf("client: could unmarshal response body: %s\n", err)
		s.log.Infof(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	return &sqlResult, nil
}

func NewStore(log *zap.SugaredLogger, elasticSearchUrl string) *Store {
	return &Store{
		log:              log,
		elasticSearchUrl: elasticSearchUrl,
	}
}
