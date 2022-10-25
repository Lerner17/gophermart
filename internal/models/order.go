package models

import "time"

type Order struct {
	UserID     int       `json:"user_id" db:"user_id"`
	Number     int64     `json:"number" db:"number"`
	Status     string    `json:"status" db:"status"`
	Accrual    int64     `json:"accrual" db:"accrual"`
	UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
}
