package main

import (
	"net/http"

	"github.com/webdevfuel/go-htmx-data-dashboard/db"
	"github.com/webdevfuel/go-htmx-data-dashboard/handler"
	"github.com/webdevfuel/go-htmx-data-dashboard/router"
	"github.com/webdevfuel/go-htmx-data-dashboard/search"
)

func main() {
	bundb := db.NewBunDB()

	meilisearchClient := search.NewMeilisearchClient()

	handler := handler.NewHandler(bundb, meilisearchClient)

	r := router.NewRouter(handler, http.Dir("./static"))

	http.ListenAndServe("localhost:3000", r)
}
