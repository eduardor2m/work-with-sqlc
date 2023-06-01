package sqlite

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
)

func GetConnection() (*sqlx.DB, error) {
	conn, err := sqlx.Connect("sqlite3", "./sqlite.db")

	if err != nil {
		return nil, err
	}

	str := "CREATE TABLE IF NOT EXISTS author (id INTEGER PRIMARY KEY, name TEXT NOT NULL, bio TEXT);"

	statement, err := conn.Prepare(str)

	if err != nil {
		return nil, err
	}

	_, err = statement.Exec()

	if err != nil {
		return nil, err
	}

	err = statement.Close()

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func CloseConnection(conn *sqlx.DB) {
	err := conn.Close()

	if err != nil {
		log.Error(err)
	}
}
