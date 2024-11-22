package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/form"
	"github.com/google/uuid"
	"github.com/meilisearch/meilisearch-go"
	"github.com/uptrace/bun"
	"github.com/webdevfuel/go-htmx-data-dashboard/data"
	"github.com/webdevfuel/go-htmx-data-dashboard/live"
	"github.com/webdevfuel/go-htmx-data-dashboard/pagination"
	"github.com/webdevfuel/go-htmx-data-dashboard/validation"
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

	err := h.bundb.NewRaw("select id, name, email, status from users where id = ?",
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

func (h *Handler) NewUserHandler(w http.ResponseWriter, r *http.Request) {
	component := view.NewUser()
	component.Render(r.Context(), w)
}

func (h *Handler) RefreshChartHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var column string
	if id == "new-users" {
		column = "new_users"
	}
	if id == "new-activations" {
		column = "new_activations"
	}

	metrics := make([]data.Metric, 0)
	err := data.GetMetric(r.Context(), h.bundb, column, &metrics)
	if err != nil {
		log.Println(err)
		return
	}

	if id == "new-users" {
		component := view.Chart(
			id,
			data.MetricsDates(metrics),
			data.NewUsersMetrics(metrics),
			true,
		)
		component.Render(r.Context(), w)
	}

	if id == "new-activations" {
		component := view.Chart(
			id,
			data.MetricsDates(metrics),
			data.NewActivationsMetrics(metrics),
			true,
		)
		component.Render(r.Context(), w)
	}

}

type NewUserFormData struct {
	Name   string `form:"name"   validate:"required"`
	Email  string `form:"email"  validate:"required"`
	Status string `form:"status" validate:"required"`
}

var decoder *form.Decoder

func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder = form.NewDecoder()

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	var data NewUserFormData

	err = decoder.Decode(&data, r.Form)
	if err != nil {
		log.Println(err)
		return
	}

	v := validation.New()

	err = v.Struct(data)
	if err != nil {
		errors := validation.Errors(data, err)
		component := view.NewUserForm(errors)
		component.Render(r.Context(), w)
	}

	id, _ := uuid.NewRandom()

	_, err = h.bundb.NewRaw(
		"insert into users (email, name, status, created_at) values (?, ?, ?, ?);",
		data.Email,
		data.Name,
		data.Status,
		time.Now().Format("2007-01-02 15:04:05"),
	).Exec(r.Context())
	if err != nil {
		log.Println("error insert into database:", err)
	}

	_, err = h.meilisearchClient.Index("users").AddDocuments([]map[string]any{
		{
			"id":     id,
			"name":   data.Name,
			"email":  data.Email,
			"status": data.Status,
		},
	})
	if err != nil {
		log.Println("error adding document to meilisearch:", err)
	}

	h.notification.Broadcast <- []byte(fmt.Sprintf(`<div id="notifications" hx-swap-oob="beforeend"><p>New user - ID: %s - Name: %s</p></div>`, id, data.Name))

	w.Header().Set("Hx-Redirect", "/users")
}

func (h *Handler) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	metrics := make([]data.Metric, 0)
	err := data.GetMetrics(r.Context(), h.bundb, &metrics)
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
