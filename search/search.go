package search

import (
	"os"

	"github.com/meilisearch/meilisearch-go"
)

func NewMeilisearchClient() meilisearch.ServiceManager {
	return meilisearch.New(os.Getenv("MEILISEARCH_HOST"))
}
