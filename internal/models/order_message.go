package models

type OrderMessage struct {
	OrderID     int    `json:"order_id"`
	OrderNumber string `json:"order_number"`
	UserID      int    `json:"user_id"`
}
