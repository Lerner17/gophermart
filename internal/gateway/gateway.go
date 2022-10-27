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

	result := &models.AccrualOrder{}

	url := fmt.Sprintf("http://%s/api/orders/%s", cfg.AccrualSystemAddress, orderNumber)

	client := request.Client{URL: url, Method: "GET"}

	resp := client.Send().Scan(&result)
	if !resp.OK() {
		// handle error
		log.Println(resp.Error())
	}

	fmt.Println(result)

	err := db.UpdateOrderState(ctx, orderID, result.Status, userID, result.Accrual)
	if err != nil {
		fmt.Println(err)
	}
}
