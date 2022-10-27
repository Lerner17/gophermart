package gateway

import (
	"context"
	"fmt"
	"log"

	"github.com/Lerner17/gophermart/internal/config"
	"github.com/Lerner17/gophermart/internal/models"

	"github.com/monaco-io/request"
)

type OrderUpdater interface {
	UpdateOrderState(context.Context, int, string, int, float64) error
}

func CalculateBonuce(db OrderUpdater, orderID int, orderNumber string, userID int) {
	ctx := context.Background()
	cfg := config.Instance

	order := models.AccrualOrder{}

	url := fmt.Sprintf("http://%s/api/orders/%s", cfg.AccrualSystemAddress, orderNumber)

	client := request.Client{URL: url, Method: "GET"}

	resp := client.Send()
	fmt.Println(resp.ScanJSON(&order))
	if !resp.OK() {
		// handle error
		log.Println(resp.Error())
	}

	// json.Unmarshal(result, &)

	err := db.UpdateOrderState(ctx, orderID, order.Status, userID, order.Accrual)
	if err != nil {
		fmt.Println(err)
	}
}
