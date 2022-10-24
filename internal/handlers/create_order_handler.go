package handlers

import (
	"context"
	"fmt"
	"io"
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

var ErrInvalidRequestFormat = &er.HTTPError{
	Code: 400,
	Msg:  "invalid request format",
}

type DBOrderCreator interface {
	CreateOrder(context.Context, models.Order) error
}

func CreateOrderHandler(db DBOrderCreator) echo.HandlerFunc {
	return func(c echo.Context) error {

		body := c.Request().Body
		data, err := io.ReadAll(body)

		if err != nil {
			return fmt.Errorf("invalid request format: %v: %w", err, ErrInvalidRequestFormat)
		}
		defer body.Close()

		token, err := c.Cookie("access-token")

		if err != nil {
			return err
		}

		userID, err := auth.GetUserIdFromToken(token.Value)
		if err != nil {
			return err
		}

		orderNumber, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil || !helpers.ValidLuhn(int(orderNumber)) {
			return fmt.Errorf("invalid order number: %v: %w", err, ErrInvalidOrderNumber)
		}

		var order = models.Order{
			UserID: userID,
			Number: orderNumber,
			Status: "NEW",
		}

		ctx := c.Request().Context()
		if err := db.CreateOrder(ctx, order); err != nil {

		}

		return nil
	}
}
