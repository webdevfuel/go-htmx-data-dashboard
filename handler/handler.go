package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/meilisearch/meilisearch-go"
	"github.com/uptrace/bun"
	"github.com/webdevfuel/go-htmx-data-dashboard/data"
	"github.com/webdevfuel/go-htmx-data-dashboard/live"
	"github.com/webdevfuel/go-htmx-data-dashboard/pagination"
	"github.com/webdevfuel/go-htmx-data-dashboard/view"
)

type Handler struct {
	bundb             *bun.DB
	meilisearchClient meilisearch.ServiceManager
	notification      *live.Notification
}

func NewHandler(
	bundb *bun.DB,
	meilisearchClient meilisearch.ServiceManager,
	notification *live.Notification,
) *Handler {
	return &Handler{
		bundb:             bundb,
		meilisearchClient: meilisearchClient,
		notification:      notification,
	}
}

func (h *Handler) UserHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var user data.User

	err := h.bundb.NewRaw("select id, name, email, status from ?.? where id = ?",
		bun.Ident("public_01_initial"),
		bun.Ident("users"),
		id,
	).Scan(r.Context(), &user)
	if err != nil {
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

	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return
	}

	searchRes, err := h.meilisearchClient.Index("users").
		SearchWithContext(r.Context(), "", &meilisearch.SearchRequest{
			Sort:        []string{sort},
			Filter:      []string{filter},
			Page:        int64(page),
			HitsPerPage: int64(25),
		})
	if err != nil {
		return
	}

	var users []data.User

	for _, hit := range searchRes.Hits {
		s, err := json.Marshal(hit)
		if err != nil {
			return
		}
		var user data.User
		err = json.Unmarshal(s, &user)
		if err != nil {
			return
		}
		users = append(users, user)
	}

	component := view.UsersTable(
		users,
		pagination.PrevPage(searchRes.Page),
		pagination.NextPage(searchRes.Page, searchRes.TotalPages),
	)
	component.Render(r.Context(), w)
}

func (h *Handler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	component := view.Users()
	component.Render(r.Context(), w)
}

func (h *Handler) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	metrics := make([]data.Metric, 0)

	err := h.bundb.NewRaw(
		"select metric_date, new_users, new_activations from ?.?;",
		bun.Ident("public_01_initial"),
		bun.Ident("metrics"),
	).Scan(r.Context(), &metrics)
	if err != nil {
		log.Println(err)
		return
	}

	component := view.Dashboard(metrics)
	component.Render(r.Context(), w)
}

func (h *Handler) Live(w http.ResponseWriter, r *http.Request) {
	c, err := live.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	client := &live.Client{Conn: c, Notification: h.notification, Send: make(chan []byte)}
	h.notification.Register <- client
	go client.Pump()
}

func (h *Handler) Notification(w http.ResponseWriter, r *http.Request) {
	h.notification.Broadcast <- []byte(`<div id="notifications" hx-swap-oob="beforeend"><p>New notification</p></div>`)
}
