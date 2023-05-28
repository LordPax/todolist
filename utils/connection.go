package utils

import (
	"database/sql"

	_ "github.com/glebarez/go-sqlite"
)

type Connection struct {
	DB *sql.DB
}

var SqliteInstance Connection

func ConnectDB(memory bool) (Connection, error) {
	filename := ":memory:"

	if !memory {
		filename = DB_FILE
	}

	db, err := sql.Open("sqlite", filename)

	if err != nil {
		return Connection{}, err
	}

	db.Exec("CREATE TABLE IF NOT EXISTS tasks (id INTEGER PRIMARY KEY, name TEXT, description TEXT, completed BOOLEAN, end_date DATETIME, begin_date DATETIME, priority INTEGER, location TEXT, label TEXT, user_id INTEGER, created_at DATETIME, updated_at DATETIME)")

	db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT, email TEXT, birthdate DATETIME, password TEXT, created_at DATETIME, updated_at DATETIME)")

	return Connection{db}, nil
}

func (c *Connection) Close() error {
	return c.DB.Close()
}

func (c *Connection) ClearDB() error {
	_, err := c.DB.Exec("DELETE FROM tasks")
	_, err = c.DB.Exec("DELETE FROM users")

	if err != nil {
		return err
	}

	return nil
}
