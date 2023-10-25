package queries

import (
	"omokogo/globals"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func InitDB() error {
	dbc, err := sqlx.Connect("sqlite3", globals.Env.DB_URL)
	if err != nil {
		return err
	}
	db = dbc
	return nil
}
