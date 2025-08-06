package database

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/silver-volt4/swapdoodle/globals"
)

var Postgres *sql.DB

func ConnectPostgres(uri string) {
	var err error

	Postgres, err = sql.Open("postgres", uri)
	if err != nil {
		globals.Logger.Critical(err.Error())
	}

	globals.Logger.Success("Connected to Postgres!")
}
