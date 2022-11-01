package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/Lerner17/gophermart/internal/config"
	er "github.com/Lerner17/gophermart/internal/errors"
	"github.com/Lerner17/gophermart/internal/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var instance *Database

type Database struct {
	cursor *sql.DB
}

func New() *Database {
	dsn := config.Instance.DatabaseDsn
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
	return instance
}

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func (db Database) CreateOrder(ctx context.Context, order models.Order) (int, error) {

	var id int

	var stmt = psql.RunWith(db.cursor).Insert("orders").SetMap(map[string]interface{}{
		"order_number": order.Number,
		"user_id":      order.UserID,
		"status":       order.Status,
	}).Suffix("RETURNING \"id\"")

	err := stmt.QueryRow().Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			err := db.checkOrder(ctx, order.Number, int64(order.UserID))
			if err != nil {
				return id, err
			}
			return id, er.ErrOrderNumberAlreadyExists
		}
		return id, fmt.Errorf("could not insert order: %v", err)
	}
	return id, nil
}

func (db Database) GetUserBalance(ctx context.Context, userID int) (models.Balance, error) {
	var balance = models.Balance{}

	query := `
		select
			sum(amount) as balance,
			(
			select
				sum(amount)
			from
				transactions t2
			where
				t2.amount < 0
				and user_id = $1) * -1  as withdrawn
		from
			transactions t
		where
			user_id = $1
	`
	fmt.Println(query)

	err := db.cursor.QueryRowContext(ctx, query, userID).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return balance, err
	}

	return balance, nil
}

func (db Database) CreateTransaction(ctx context.Context, userID int, orderNum string, amount float64) error {

	orderNumber, err := strconv.ParseInt(string(orderNum), 10, 64)

	if err != nil {
		return err
	}

	orderID, err := db.getOrderID(userID, orderNumber)

	if err != nil {
		return err
	}

	var stmt = psql.RunWith(db.cursor).Insert("transactions").SetMap(map[string]interface{}{
		"user_id":  userID,
		"order_id": orderID,
		"amount":   amount * -1,
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return fmt.Errorf("could not insert transaction: %v", err)
	}

	return nil
}
