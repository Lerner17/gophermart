package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	er "github.com/Lerner17/gophermart/internal/errors"
	"github.com/Lerner17/gophermart/internal/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
)

func (db Database) CreateOrderWithWithdraws(ctx context.Context, userID int, o models.Order) error {

	var totalBalance sql.NullFloat64

	tx, err := db.cursor.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not start transaction: %v", err)
	}
	defer tx.Rollback() // Because we are run Rollback after Commit (at the end of function) it would be OK

	query := psql.Select("sum(amount)").From("transactions").Where(sq.Eq{
		"user_id": userID,
	}).RunWith(tx).PlaceholderFormat(sq.Dollar)

	if err := query.QueryRow().Scan(&totalBalance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return er.ErrCannotFindTransactions
		}
		return err
	}

	if !totalBalance.Valid || totalBalance.Float64 < float64(o.Accrual.Float64) {
		return er.ErrBalanceTooLow
	}

	var oid int64
	stmt, err := tx.PrepareContext(ctx, `
		insert into orders(
			user_id
			,order_number
			,status
			,processed_at
		) values($1, $2, $3, $4) returning id`)

	if err != nil {
		return err
	}

	err = stmt.QueryRowContext(ctx, userID, o.Number, "PROCESSED", time.Now()).Scan(&oid)
	if err != nil {
		return fmt.Errorf("cannot insert order: %v", err)
	}

	_, err = tx.ExecContext(ctx, "insert into transactions (user_id, order_id, amount) values ($1, $2, $3)", userID, oid, -1*o.Accrual.Float64)

	if err != nil {
		return fmt.Errorf("cannot insert transaction: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("cannot commit transaction: %v", err)
	}

	return nil
}

func (db Database) GetNewOrders(ctx context.Context) ([]models.Order, error) {
	orders := make([]models.Order, 0)

	query := psql.Select("id", "order_number", "status").From("orders").
		Where(sq.Eq{"status": "NEW"}).
		RunWith(db.cursor).
		PlaceholderFormat(sq.Dollar)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return orders, err
	}

	if rows.Err() != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return orders, er.ErrOrdersNotFound
		}
	}

	defer rows.Close()
	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.ID, &order.Number, &order.Status)
		if err != nil {
			return orders, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (db Database) GetOrders(ctx context.Context, userID int) ([]models.Order, error) {

	orders := make([]models.Order, 0)
	query := psql.Select("order_number", "status", "amount as accrual", "uploaded_at").
		From("orders o").
		Where(sq.Eq{"o.user_id": userID}).
		LeftJoin("transactions t on o.id = t.order_id").
		RunWith(db.cursor).
		PlaceholderFormat(sq.Dollar)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return orders, err
	}

	if rows.Err() != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return orders, er.ErrOrdersNotFound
		}
	}
	defer rows.Close()
	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return orders, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (db Database) UpdateOrderState(ctx context.Context, orderID int, orderStatus string, userID int, amount float64) error {
	tx, err := db.cursor.Begin()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "update orders o set status=$1 where o.id = $2", orderStatus, orderID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = tx.ExecContext(ctx, "insert into transactions (user_id, order_id, amount) values ($1, $2, $3)", userID, orderID, amount)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err = tx.Commit(); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (db Database) checkOrder(ctx context.Context, orderNumber string, userID int64) error {
	var uid int64

	query := psql.
		Select("user_id").
		From("orders").
		Where(sq.Eq{"order_number": orderNumber}).
		RunWith(db.cursor).
		PlaceholderFormat(sq.Dollar)

	if err := query.QueryRow().Scan(&uid); err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return er.ErrCannotFindOrderByNumber
		}
		return err
	}

	if userID == uid {
		return er.ErrOrderWasCreatedBySelf
	} else {
		return er.ErrOrderWasCreatedByAnotherUser
	}
}

func (db Database) getOrderID(userID int, orderNumber int64) (int, error) {
	var id int
	query := psql.Select("id").From("orders").Where(sq.Eq{
		"order_number": orderNumber,
		"user_id":      userID,
	}).RunWith(db.cursor).PlaceholderFormat(sq.Dollar)
	if err := query.QueryRow().Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return 0, er.ErrCannotFindOrderByNumber
		}
		return id, err
	}
	return id, nil
}

func (db Database) GetWithdraws(ctx context.Context, userID int) ([]models.Withdraw, error) {
	withdraw := make([]models.Withdraw, 0)
	query := psql.Select("order_number", "processed_at", "t.amount * -1 as amount").
		From("orders o").
		LeftJoin("transactions t on o.id = t.order_id").
		Where(sq.Eq{"o.user_id": userID}).
		Where(sq.Lt{"t.amount": 0}).
		RunWith(db.cursor).
		PlaceholderFormat(sq.Dollar)
	fmt.Println(query.ToSql())
	rows, err := query.QueryContext(ctx)

	if err != nil {
		return withdraw, err
	}

	if rows.Err() == nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return withdraw, er.ErrOrdersNotFound
		}
	}

	defer rows.Close()
	for rows.Next() {
		var w models.Withdraw
		fmt.Println(w)
		err = rows.Scan(&w.Number, &w.ProcessedAt, &w.Sum)
		if err != nil {
			return withdraw, err
		}
		fmt.Println(w)
		withdraw = append(withdraw, w)
	}
	return withdraw, nil
}
