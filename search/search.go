package search

import (
	"context"
	"encoding/json"
	"os"

	"github.com/meilisearch/meilisearch-go"
)

var meilisearchClient = meilisearch.New(os.Getenv("MEILISEARCH_HOST"))

func Search[T any](ctx context.Context, values []T, sort, filter string) ([]T, error) {
	searchRes, err := meilisearchClient.Index("users").
		SearchWithContext(ctx, "", &meilisearch.SearchRequest{
			Sort:   []string{sort},
			Filter: []string{filter},
		})
	if err != nil {
		return nil, err
	}

	for _, hit := range searchRes.Hits {
		s, err := json.Marshal(hit)
		if err != nil {
			return nil, err
		}
		var value T
		err = json.Unmarshal(s, &value)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil
}
