package search

import (
	"context"
	"os"

	"github.com/meilisearch/meilisearch-go"
)

func NewMeilisearchClient() meilisearch.ServiceManager {
	return meilisearch.New(os.Getenv("MEILISEARCH_HOST"))
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
