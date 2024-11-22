package data

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

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

func NewUsersMetrics(metrics []Metric) []int {
	newUsers := make([]int, 0)
	for _, metric := range metrics {
		newUsers = append(newUsers, metric.NewUsers)
	}
	return newUsers
}

func NewActivationsMetrics(metrics []Metric) []int {
	newActivations := make([]int, 0)
	for _, metric := range metrics {
		newActivations = append(newActivations, metric.NewActivations)
	}
	return newActivations
}

func MetricsDates(metrics []Metric) []string {
	dates := make([]string, 0)
	for _, metric := range metrics {
		d, _ := time.Parse("2006-01-02", metric.MetricDate)
		dates = append(dates, d.Format("02 Jan"))
	}
	return dates
}

func GetMetric(ctx context.Context, bundb *bun.DB, column string, dest *[]Metric) error {
	err := bundb.NewRaw(
		"select metric_date, ? from metrics;",
		bun.Ident(column),
	).Scan(ctx, dest)
	if err != nil {
		return nil
	}
	return err
}

func GetMetrics(ctx context.Context, bundb *bun.DB, dest *[]Metric) error {
	err := bundb.NewRaw(
		"select metric_date, new_users, new_activations from metrics;",
	).Scan(ctx, dest)
	if err != nil {
		return err
	}
	return nil
}
