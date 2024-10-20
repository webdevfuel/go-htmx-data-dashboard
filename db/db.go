package db

import (
	"database/sql"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func NewBunDB() *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(os.Getenv("DATABASE_URL"))))
	bundb := bun.NewDB(sqldb, pgdialect.New())
	bundb.AddQueryHook(bundebug.NewQueryHook())
	return bundb
}
