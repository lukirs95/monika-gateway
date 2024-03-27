package sqlite

import (
	"database/sql"

	sdk "github.com/lukirs95/monika-gosdk/pkg/types"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteRepository struct {
	db            *sql.DB
	adminPassword string
}

func NewSQLiteRepository(db *sql.DB, adminPassword string) *SQLiteRepository {
	return &SQLiteRepository{
		db:            db,
		adminPassword: adminPassword,
	}
}

// Migrate validates database and creates required tables if not present
func (s *SQLiteRepository) Migrate() error {
	// check if user table is present
	userRow := s.db.QueryRow(EXISTS_TABLE_USERS)
	var userTable string
	if err := userRow.Scan(&userTable); err != nil {
		if err == sql.ErrNoRows {
			// create table users in database
			if _, err := s.db.Exec(CREATE_TABLE_USERS); err != nil {
				return err
			}
			if _, err := s.db.Exec(CREATE_USER, "admin", s.adminPassword, sdk.UserRole_ADMIN); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// check if group table is present
	groupsRow := s.db.QueryRow(EXISTS_TABLE_GROUPS)
	var groupsTable string
	if err := groupsRow.Scan(&groupsTable); err != nil {
		if err == sql.ErrNoRows {
			// create table groups in database
			if _, err := s.db.Exec(CREATE_TABLE_GROUPS); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// check if modules table is present
	modulesRow := s.db.QueryRow(EXISTS_TABLE_MEMBERS)
	var modulesTable string
	if err := modulesRow.Scan(&modulesTable); err != nil {
		if err == sql.ErrNoRows {
			// create table modules in database
			if _, err := s.db.Exec(CREATE_TABLE_MEMBERS); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}
