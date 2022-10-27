package handlers

import (
	"context"

	"github.com/labstack/echo/v4"
)

type WithdrawGetter interface {
	GetWithdraws(context.Context, int) error
}

func GetWithdrawHandler(db WithdrawGetter) echo.HandlerFunc {
	return func(c echo.Context) error {

		return nil
	}
}
