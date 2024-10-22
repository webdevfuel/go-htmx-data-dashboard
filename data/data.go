package data

type User struct {
	ID     int    `json:"id"     db:"id"`
	Name   string `json:"name"   db:"name"`
	Email  string `json:"email"  db:"email"`
	Status string `json:"status" db:"status"`
}

type Metric struct {
	MetricDate     string `db:"metric_date"`
	NewUsers       int    `db:"new_users"`
	NewActivations int    `db:"new_activations"`
}
