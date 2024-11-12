package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	randmath "math/rand/v2"

	"github.com/meilisearch/meilisearch-go"
	"github.com/uptrace/bun"
	"github.com/webdevfuel/go-htmx-data-dashboard/db"
	"github.com/webdevfuel/go-htmx-data-dashboard/search"
)

var statuses []string = []string{
	"active",
	"blocked",
	"pending",
	"archived",
}

var dates []string = []string{
	"2022-10-22",
	"2022-10-21",
	"2022-10-20",
	"2022-10-19",
	"2022-10-18",
	"2022-10-17",
	"2022-10-16",
	"2022-10-15",
}

func main() {
	bundb := db.NewBunDB()

	_, err := bundb.NewRaw("delete from ?.?;",
		bun.Ident("public_01_initial"),
		bun.Ident("users"),
	).Exec(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	_, err = bundb.NewRaw("delete from ?.?;",
		bun.Ident("public_01_initial"),
		bun.Ident("metrics"),
	).Exec(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	meilisearchClient := search.NewMeilisearchClient()

	_, err = meilisearchClient.Index("users").
		DeleteAllDocuments()
	if err != nil {
		log.Fatal(err)
	}

	documents := make([]map[string]any, 0)

	for i := range 5000 {
		var id int
		name := randomString()
		status := statuses[randmath.IntN(4)]
		date := dates[randmath.IntN(8)]
		createdAt := fmt.Sprintf("%s 00:00:00", date)
		var activatedAt *string
		if status == "active" {
			s := fmt.Sprintf("%s 00:00:00", date)
			activatedAt = &s
		}
		email := fmt.Sprintf("%s@example.com", name)
		err := bundb.NewRaw(
			"insert into ?.? (email, name, status, created_at, activated_at) values (?, ?, ?, ?, ?) returning id;",
			bun.Ident("public_01_initial"),
			bun.Ident("users"),
			email,
			name,
			status,
			createdAt,
			activatedAt,
		).Scan(context.Background(), &id)
		if err != nil {
			log.Fatal(err)
		}
		documents = append(documents, map[string]any{
			"id":     id,
			"name":   name,
			"email":  email,
			"status": status,
		})
		if i%1000 == 0 {
			_, err := meilisearchClient.Index("users").
				AddDocuments(documents)
			if err != nil {
				log.Fatal(err)
			}
			documents = make([]map[string]any, 0)
		}
	}

	for _, date := range dates {
		public01Initial := bun.Ident("public_01_initial")
		_, err := bundb.NewRaw(`
			INSERT INTO ?.? (metric_date, new_users, new_activations)
			SELECT
				to_date(?, 'YYYY-MM-DD') AS metric_date,
				(SELECT COUNT(*) FROM ?.? WHERE created_at::date = to_date(?, 'YYYY-MM-DD')) AS new_users,
				(SELECT COUNT(*) FROM ?.? WHERE activated_at::date = to_date(?, 'YYYY-MM-DD')) AS new_activations
			ON CONFLICT (metric_date)
				DO UPDATE SET
					new_users = EXCLUDED.new_users,
					new_activations = EXCLUDED.new_activations;
		`, public01Initial, bun.Ident("metrics"), date, public01Initial, bun.Ident("users"), date, public01Initial, bun.Ident("users"), date).Exec(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = meilisearchClient.CreateIndex(&meilisearch.IndexConfig{
		Uid:        "users",
		PrimaryKey: "id",
	})
	if err != nil {

	}

	sortableAttributes := []string{"name", "email", "status"}
	_, err = meilisearchClient.Index("users").
		UpdateSortableAttributes(&sortableAttributes)
	if err != nil {

	}

	filterableAttributes := []string{"status"}
	_, err = meilisearchClient.Index("users").
		UpdateFilterableAttributes(&filterableAttributes)
	if err != nil {

	}
}

func randomString() string {
	b := make([]byte, 12)
	_, err := rand.Read(b)
	if err != nil {

	}
	return base64.StdEncoding.EncodeToString(b)
}
