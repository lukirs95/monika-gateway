package db_test

import (
	"database/sql"
	"testing"

	monikadb "github.com/lukirs95/monika-gateway/internal/monika-db"
	_ "github.com/mattn/go-sqlite3"
)

func TestDB(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	var version string
	err = db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(version)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	_, close := monikadb.NewSQLiteDatabase(":memory:", "admin")
	if err := close(); err != nil {
		t.Fatal(err)
	}
}
