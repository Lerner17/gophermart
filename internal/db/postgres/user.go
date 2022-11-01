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
)

func (db Database) LoginUser(username, password string) (int, error) {
	var id int
	var p string
	query := psql.Select("id", "password").From("users").Where(sq.Eq{"username": username}).RunWith(db.cursor).PlaceholderFormat(sq.Dollar)

	if err := query.QueryRow().Scan(&id, &p); err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return 0, er.ErrInvalidLoginOrPassword
		}
		return 0, err
	}
	fmt.Println(id)
	if verefyPassword := helpers.ComparePasswords(p, []byte(password)); !verefyPassword {
		return 0, er.ErrInvalidLoginOrPassword
	}

	return id, nil
}

func (db Database) RegisterUser(ctx context.Context, username, password string) (int, error) {
	var id int
	hashedPassword, err := helpers.HashAndSalt([]byte(password))
	if err != nil {
		return id, fmt.Errorf("could not hash password: %v", err)
	}
	var stmt = psql.RunWith(db.cursor).Insert("users").SetMap(map[string]interface{}{
		"username": username,
		"password": hashedPassword,
	}).Suffix("RETURNING \"id\"")

	err = stmt.QueryRowContext(ctx).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return id, er.ErrUserNameAlreadyExists
		}
		return id, fmt.Errorf("could not insert user: %v", err)
	}
	return id, nil
}
