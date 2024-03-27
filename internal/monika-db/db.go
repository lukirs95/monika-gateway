package monikadb

import (
	"database/sql"

	monikadb "github.com/lukirs95/monika-gateway/internal/monika-db/dbrepo"
	"github.com/lukirs95/monika-gateway/internal/monika-db/sqlite"
)

func NewSQLiteDatabase(path string, adminPassword string) (monikadb.DatabaseRepo, func() error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}

	return sqlite.NewSQLiteRepository(db, adminPassword), db.Close
}
