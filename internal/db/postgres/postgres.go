package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/Lerner17/gophermart/internal/config"
	er "github.com/Lerner17/gophermart/internal/errors"
	"github.com/Lerner17/gophermart/internal/helpers"
	"github.com/Lerner17/gophermart/internal/models"
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

func (db Database) CreateOrder(ctx context.Context, order models.Order) error {

	var stmt = psql.RunWith(db.cursor).Insert("orders").SetMap(map[string]interface{}{
		"number":  order.Number,
		"user_id": order.UserID,
		"status":  order.Status,
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			err := db.checkOrder(ctx, order.Number, int64(order.UserID))
			if err != nil {
				return err
			}
			return er.OrderNumberAlreadyExists
		}
		fmt.Println(err)
		return fmt.Errorf("could not insert order: %v", err)
	}
	return nil
}

func (db Database) LoginUser(username, password string) (int, error) {
	var id int
	var p string
	query := psql.Select("id", "password").From("users").Where(sq.Eq{"username": username}).RunWith(db.cursor).PlaceholderFormat(sq.Dollar)

	if err := query.QueryRow().Scan(&id, &p); err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return 0, er.InvalidLoginOrPassword
		}
		return 0, err
	}
	fmt.Println(id)
	if verefyPassword := helpers.ComparePasswords(p, []byte(password)); !verefyPassword {
		return 0, er.InvalidLoginOrPassword
	}

	return id, nil
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

func (db Database) GetUserBalance(ctx context.Context, userID int) (models.Balance, error) {
	var balance = models.Balance{}

	_ = psql.Select("sum(amount), ").From("transactions").Where(sq.Eq{"user_id": userID}).RunWith(db.cursor).PlaceholderFormat(sq.Dollar)

	return balance, nil
}

func (db Database) getOrderID(userID int, orderNumber int64) (int, error) {
	var id int
	query := psql.Select("id").From("orders").Where(sq.Eq{
		"number":  orderNumber,
		"user_id": userID,
	}).RunWith(db.cursor).PlaceholderFormat(sq.Dollar)
	if err := query.QueryRow().Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return 0, er.CannotFindOrderByNumber
		}
		return id, err
	}
	return id, nil
}

func (db Database) checkUserBalance(userID int, amount int) error {

	var totalBalance int

	// select sum(amount) from transactions t where user_id = 1
	query := psql.Select("sum(amount)").From("transactions").Where(sq.Eq{
		"user_id": userID,
	}).RunWith(db.cursor).PlaceholderFormat(sq.Dollar)

	if err := query.QueryRow().Scan(&totalBalance); err != nil {
		// if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
		// 	return , er.CannotFindOrderByNumber TODO:
		// }
		return err
	}
	if totalBalance < amount {
		return errors.New("user balance too low")
	}

	return nil
}

func (db Database) CreateTransaction(ctx context.Context, userID int, orderNum string, amount int) error {

	orderNumber, err := strconv.ParseInt(string(orderNum), 10, 64)

	if err = db.checkUserBalance(userID, amount); err != nil {
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

func (db Database) GetOrders(ctx context.Context, userID int) ([]models.Order, error) {

	orders := make([]models.Order, 0)
	query := psql.Select("number", "status", "accrual", "uploaded_at").From("orders").Where(sq.Eq{"user_id": userID}).RunWith(db.cursor).PlaceholderFormat(sq.Dollar)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return orders, err
	}

	if rows.Err() == nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return orders, er.OrdersNotFound
		}
	}
	defer rows.Close()
	for rows.Next() {
		var order models.Order
		// var orderTime string
		err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return orders, err
		}

		// order.UploadedAt, err = time.Parse("2006-01-02T15:04:05.000Z", orderTime)
		if err != nil {
			return orders, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (db Database) checkOrder(ctx context.Context, orderNumber int64, userID int64) error {
	var uid int64

	query := psql.Select("user_id").From("orders").Where(sq.Eq{"number": orderNumber}).RunWith(db.cursor).PlaceholderFormat(sq.Dollar)

	if err := query.QueryRow().Scan(&uid); err != nil {
		// if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
		// 	return 0, er.InvalidLoginOrPassword
		// }
		return err
	}

	if userID == uid {
		return er.OrderWasCreatedBySelf
	} else {
		return er.OrderWasCreatedByAnotherUser
	}
}
