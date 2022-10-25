package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Lerner17/gophermart/internal/auth"
	er "github.com/Lerner17/gophermart/internal/errors"
	"github.com/Lerner17/gophermart/internal/helpers"
	"github.com/Lerner17/gophermart/internal/models"
	"github.com/labstack/echo/v4"
)

var ErrInvalidRequestFormat = &er.HTTPError{
	Code: 400,
	Msg:  "invalid request format",
}

var ErrOrderAlreadyExists = &er.HTTPError{
	Code: 200,
	Msg:  "order already exists",
}

var ErrOrderAlreadyExistsByAnotherUser = &er.HTTPError{
	Code: 409,
	Msg:  "order already exists by another user",
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
			if err != nil {
				if errors.Is(err, er.OrderWasCreatedByAnotherUser) {
					return fmt.Errorf("conflict: %v: %w", err, ErrOrderAlreadyExistsByAnotherUser)
				}
				if errors.Is(err, er.OrderWasCreatedBySelf) {
					return fmt.Errorf("already exists: %v: %w", err, ErrOrderAlreadyExists)
				}
				return err
			}
			return fmt.Errorf("confilct order number: %v: %w", err, ErrInvalidOrderNumber)
		}
		return c.String(http.StatusAccepted, "")
	}
}
