package search

import (
	"context"
	"encoding/json"
	"os"

	"github.com/meilisearch/meilisearch-go"
)

func NewMeilisearchClient() meilisearch.ServiceManager {
	return meilisearch.New(os.Getenv("MEILISEARCH_HOST"))
}

func Search[T any](
	ctx context.Context,
	client meilisearch.ServiceManager,
	values []T,
	sort, filter string,
) ([]T, error) {
	searchRes, err := client.Index("users").
		SearchWithContext(ctx, "", &meilisearch.SearchRequest{
			Sort:   []string{sort},
			Filter: []string{filter},
			Limit:  100,
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

func InsertDocuments(
	ctx context.Context,
	client meilisearch.ServiceManager,
	documents any,
) error {
	_, err := client.Index("users").
		AddDocuments(documents)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAllDocuments(
	ctx context.Context,
	client meilisearch.ServiceManager,
) error {
	_, err := client.Index("users").
		DeleteAllDocuments()
	if err != nil {
		return err
	}
	return nil
}
