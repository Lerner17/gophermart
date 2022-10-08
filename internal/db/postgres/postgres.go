package postgres

import (
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var instance *Database

type Database struct {
	cursor *sql.DB
}

func New() *Database {
	return instance
}

func init() {}

func (db *Database) RegisterUser(username, password string) error {
	return nil
}
