package consumer

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Lerner17/gophermart/internal/config"
	"github.com/Lerner17/gophermart/internal/models"
	"github.com/Lerner17/gophermart/internal/queue"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
)

type OrderUpdater interface {
	UpdateOrderState(context.Context, int, string, int, float64) error
}

func ProcessOrderBounce(logger echo.Logger, db OrderUpdater) {

	for {
		msg, err := queue.GetNextOrderMessage()
		if err != nil {
			if errors.Is(err, queue.ErrQueueClosed) {
				break
			}
			logger.Error("Error occured while processing message: %v", err)
		}

		ctx := context.Background()
		cfg := config.Instance

		client := resty.New()
		order := models.AccrualOrder{}

		resp, err := client.
			R().
			SetResult(&order).
			EnableTrace().
			SetContext(ctx).
			SetPathParams(map[string]string{"orderNumber": msg.Number}).
			Get(cfg.AccrualSystemAddress + "/api/orders/{orderNumber}")
		if err != nil {
			logger.Error("Bounce service is unavailable: %v", err)
			queue.PushOrderMessage(msg) // Error occured, push message back to queue
			time.Sleep(10 * time.Second)
			continue
		}

		if resp.StatusCode() == http.StatusInternalServerError {
			logger.Error("accrual system return 500 status code. retry imediatly")
			queue.PushOrderMessage(msg)
			time.Sleep(10 * time.Second)
			continue
		}

		if resp.StatusCode() == http.StatusTooManyRequests {
			logger.Error("too many requests. sleep")
			queue.PushOrderMessage(msg)
			time.Sleep(10 * time.Second)
			continue
		}

		err = db.UpdateOrderState(ctx, msg.ID, order.Status, msg.UserID, order.Accrual)
		if err != nil {
			logger.Error("Could not update order: %v", err)
			queue.PushOrderMessage(msg)
			time.Sleep(10 * time.Second)
			continue
		}
	}
}
