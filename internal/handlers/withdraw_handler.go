package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Lerner17/gophermart/internal/auth"
	er "github.com/Lerner17/gophermart/internal/errors"
	"github.com/Lerner17/gophermart/internal/helpers"
	"github.com/Lerner17/gophermart/internal/models"
	"github.com/labstack/echo/v4"
)

var ErrInvalidOrderNumber = &er.HTTPError{
	Code: 422,
	Msg:  "incorrect order number",
}

var ErrBalanceTooLow = &er.HTTPError{
	Code: 402,
	Msg:  "balance too low",
}

type WithdrawWriter interface {
	CreateTransaction(context.Context, int, string, float64) error
	CreateOrderWithWithdraws(context.Context, int, models.Order) error
}

type withdrawBody struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func WithdrawHandler(db WithdrawWriter) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := c.Cookie("access-token")

		if err != nil {
			return err
		}

		userID, err := auth.GetUserIDFromToken(token.Value)
		if err != nil {
			return err
		}

		body := new(withdrawBody)

		if err := c.Bind(&body); err != nil {
			return fmt.Errorf("could not bind body: %v: %w", err, ErrBindBody)
		}

		order := models.Order{
			UserID:      userID,
			Status:      "PROCESSED",
			Number:      body.Order,
			ProcessedAt: time.Now(),
			Accrual: sql.NullFloat64{
				Valid:   true,
				Float64: body.Sum,
			},
		}

		orderNumber, err := strconv.ParseInt(string(order.Number), 10, 64)
		fmt.Println(orderNumber)
		if err != nil || !helpers.ValidLuhn(int(orderNumber)) {
			fmt.Println(err)
			return fmt.Errorf("invalid order number: %v: %w", err, ErrInvalidOrderNumber)
		}
		fmt.Println(order)
		err = db.CreateOrderWithWithdraws(c.Request().Context(), userID, order)

		if err != nil {
			if errors.Is(err, er.ErrBalanceTooLow) {
				return fmt.Errorf("user balance too low: %v: %w", err, ErrBalanceTooLow)
			}
			return fmt.Errorf("cannot create orders with withdraw: %v", err)
		}

		return nil
	}
}
