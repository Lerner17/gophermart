package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Lerner17/gophermart/internal/helpers"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var instance *Database

type Database struct {
	cursor *sql.DB
}

func New() *Database {
	return instance
}

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func init() {
	dsn := "postgres://shroten:shroten@localhost:5432/shroten"
	if dsn == "" {
		panic("Cannot connect to database")
	}
	cursor, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	instance = &Database{
		cursor: cursor,
	}
}

func (db Database) RegisterUser(ctx context.Context, username, password string) error {
	hashedPassword, err := helpers.HashAndSalt([]byte(password))
	if err != nil {
		return fmt.Errorf("could not hash password: %v", err)
	}
	var stmt = psql.RunWith(db.cursor).Insert("users").SetMap(map[string]interface{}{
		"username": username,
		"password": hashedPassword,
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return fmt.Errorf("could not insert user: ", err)
	}
	return nil
}
