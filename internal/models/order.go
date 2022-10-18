package models

type Order struct {
	UserID int    `json:"user_id"`
	Number int64  `json:"number"`
	Status string `json:"status"`
}
