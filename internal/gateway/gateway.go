package gateway

import (
	"context"
	"fmt"

	"github.com/Lerner17/gophermart/internal/config"
	"github.com/Lerner17/gophermart/internal/models"

	"github.com/go-resty/resty/v2"
)

type OrderUpdater interface {
	UpdateOrderState(context.Context, int, string, int, float64) error
}

func CalculateBonuce(db OrderUpdater, orderID int, orderNumber string, userID int) {
	ctx := context.Background()
	cfg := config.Instance

	client := resty.New()
	order := models.AccrualOrder{}

	response, err := client.R().SetResult(&order).EnableTrace().
		SetContext(ctx).SetPathParams(map[string]string{"orderNumber": orderNumber}).
		Get(cfg.AccrualSystemAddress + "/api/orders/{orderNumber}")

	if err != nil {
		fmt.Print(err)
	}

	fmt.Println(response)

	err = db.UpdateOrderState(ctx, orderID, order.Status, userID, order.Accrual)
	if err != nil {
		fmt.Println(err)
	}
}
