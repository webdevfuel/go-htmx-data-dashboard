package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/webdevfuel/go-htmx-data-dashboard/data"
	"github.com/webdevfuel/go-htmx-data-dashboard/search"
	"github.com/webdevfuel/go-htmx-data-dashboard/view"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		component := view.Dashboard()
		component.Render(r.Context(), w)
	})
	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		component := view.Users()
		component.Render(r.Context(), w)
	})
	r.Get("/users-table", func(w http.ResponseWriter, r *http.Request) {
		sort := r.URL.Query().Get("sort")
		if sort == "" {
			sort = "name:asc"
		}

		filter := r.URL.Query().Get("filter")
		if filter == "" {
			filter = ""
		}

		var users []data.User
		users, err := search.Search(r.Context(), users, sort, filter)
		if err != nil {
			return
		}

		component := view.UsersTable(users)
		component.Render(r.Context(), w)
	})
	http.ListenAndServe("localhost:3000", r)
}
