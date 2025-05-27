package database

import (
	"database/sql"
	"os"

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

func InitNotificationTable() {
	// In our testing, having the notification ID set to low numbers caused some weird issues.
	// Bumping it up seems to have fixed them, so let's start with a high number.
	// Since notifications are sequential, this can also remedy some issues with the order not matching up (from the time when SD still worked)
	_, err := Postgres.Exec(`CREATE SEQUENCE IF NOT EXISTS datastore.notification_id_seq
		INCREMENT 1
		START 100000000
		CACHE 1`,
	)
	if err != nil {
		globals.Logger.Critical(err.Error())
		os.Exit(0)
	}

	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS datastore.notifications (
		notification_id BIGINT NOT NULL DEFAULT nextval('datastore.notification_id_seq') PRIMARY KEY,
		recipient_pid INT,
		data_id BIGINT REFERENCES datastore.objects (data_id)
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		os.Exit(0)
	}

	globals.Logger.Success("Postgres tables created")
}
