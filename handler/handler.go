package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/meilisearch/meilisearch-go"
	"github.com/uptrace/bun"
	"github.com/webdevfuel/go-htmx-data-dashboard/data"
	"github.com/webdevfuel/go-htmx-data-dashboard/search"
	"github.com/webdevfuel/go-htmx-data-dashboard/view"
)

type Handler struct {
	bundb             *bun.DB
	meilisearchClient meilisearch.ServiceManager
}

func NewHandler(bundb *bun.DB, meilisearchClient meilisearch.ServiceManager) *Handler {
	return &Handler{
		bundb:             bundb,
		meilisearchClient: meilisearchClient,
	}
}

func (h *Handler) UserHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var user data.User

	err := h.bundb.NewSelect().
		Model(&user).
		Column("id", "name", "email", "status").
		Where("id = ?", id).
		Scan(r.Context())
	if err != nil {
		log.Println(err)
		return
	}

	component := view.User(user)
	component.Render(r.Context(), w)
}

func (h *Handler) UsersTableHandler(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "name:asc"
	}

	filter := r.URL.Query().Get("filter")
	if filter == "" {
		filter = ""
	}

	var users []data.User
	users, err := search.Search(r.Context(), h.meilisearchClient, users, sort, filter)
	if err != nil {
		log.Println(err)
		return
	}

	component := view.UsersTable(users)
	component.Render(r.Context(), w)
}

func (h *Handler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	component := view.Users()
	component.Render(r.Context(), w)
}

func (h *Handler) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	component := view.Dashboard()
	component.Render(r.Context(), w)
}
