package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	randmath "math/rand/v2"

	"github.com/webdevfuel/go-htmx-data-dashboard/db"
	"github.com/webdevfuel/go-htmx-data-dashboard/search"
)

var statuses []string = []string{
	"active",
	"blocked",
	"pending",
	"archived",
}

func main() {
	bundb := db.NewBunDB()

	_, err := bundb.Exec("delete from users;")
	if err != nil {
		log.Fatal(err)
	}

	meilisearchClient := search.NewMeilisearchClient()

	err = search.DeleteAllDocuments(context.Background(), meilisearchClient)
	if err != nil {
		log.Fatal(err)
	}

	documents := make([]map[string]any, 0)

	for i := range 100000 {
		var id int
		s := randomString()
		email := fmt.Sprintf("%s@example.com", s)
		err := bundb.QueryRow(
			"insert into users (email, name, status) values (?, ?, ?) returning id;",
			email,
			s,
			"active",
		).Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		documents = append(documents, map[string]any{
			"id":     id,
			"name":   s,
			"email":  email,
			"status": statuses[randmath.IntN(4)],
		})
		if i%1000 == 0 {
			err = search.InsertDocuments(context.Background(), meilisearchClient, documents)
			if err != nil {
				log.Fatal(err)
			}
			documents = make([]map[string]any, 0)
		}
	}
}

func randomString() string {
	b := make([]byte, 12)
	_, err := rand.Read(b)
	if err != nil {

	}
	return base64.StdEncoding.EncodeToString(b)
}
