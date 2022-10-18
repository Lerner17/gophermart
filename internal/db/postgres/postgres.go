package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	er "github.com/Lerner17/gophermart/internal/errors"
	"github.com/Lerner17/gophermart/internal/helpers"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
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

func (db Database) LoginUser(username, password string) error {
	var p string
	fmt.Println(username)
	query := psql.Select("password").From("users").Where(sq.Eq{"username": username}).RunWith(db.cursor).PlaceholderFormat(sq.Dollar)

	if err := query.QueryRow().Scan(&p); err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return er.InvalidLoginOrPassword
		}
		return err
	}

	if verefyPassword := helpers.ComparePasswords(p, []byte(password)); !verefyPassword {
		return er.InvalidLoginOrPassword
	}

	return nil
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return er.UserNameAlreadyExists
		}
		return fmt.Errorf("could not insert user: %v", err)
	}
	return nil
}
