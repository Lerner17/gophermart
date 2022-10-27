package handlers

import (
	"context"
	"fmt"
	"strconv"

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
	CreateTransaction(context.Context, int, string, int) error
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

		withderaw := &models.Withdraw{}

		if err := c.Bind(withderaw); err != nil {
			return fmt.Errorf("could not bind body: %v: %w", err, ErrBindBody)
		}

		orderNumber, err := strconv.ParseInt(string(withderaw.Order), 10, 64)
		fmt.Println(orderNumber)
		if err != nil || !helpers.ValidLuhn(int(orderNumber)) {
			fmt.Println(err)
			return fmt.Errorf("invalid order number: %v: %w", err, ErrInvalidOrderNumber)
		}
		fmt.Println(withderaw)
		if err := db.CreateTransaction(c.Request().Context(), userID, withderaw.Order, withderaw.Sum); err != nil {
			fmt.Println(err)
			return fmt.Errorf("could not create withdreaw: %v: %w", err, ErrBalanceTooLow)
		}
		return nil
	}
}
