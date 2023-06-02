package sqlite

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"

	_ "github.com/lib/pq"
)

func GetConnection() (*sqlx.DB, error) {
	connStr := "postgres://root:root@localhost/mydb?sslmode=disable"
	conn, err := sqlx.Connect("postgres", connStr)

	if err != nil {
		return nil, err
	}

	str := "CREATE TABLE IF NOT EXISTS author (id SERIAL PRIMARY KEY, name TEXT NOT NULL, bio TEXT);"

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
