package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Order struct {
	UserID       int             `json:"user_id,omitempty"`
	Number       string          `json:"number"`
	Status       string          `json:"status,omitempty"`
	UploadedAt   time.Time       `json:"uploaded_at"`
	Accrual      sql.NullFloat64 `json:"accrual"`
	Processed_at time.Time       `json:"processed_at"`
}

func (o Order) MarshalJSON() ([]byte, error) {
	if !o.Accrual.Valid {
		return json.Marshal(struct {
			UserID     int       `json:"user_id,omitempty"`
			Number     string    `json:"number" db:"number"`
			Status     string    `json:"status" db:"status"`
			UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
		}{o.UserID, o.Number, o.Status, o.UploadedAt})
	}
	return json.Marshal(struct {
		UserID     int       `json:"user_id,omitempty"`
		Number     string    `json:"number" db:"number"`
		Status     string    `json:"status" db:"status"`
		UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
		Accrual    float64   `json:"accrual" db:"accrual"`
	}{o.UserID, o.Number, o.Status, o.UploadedAt, o.Accrual.Float64})
}

type AccrualOrder struct {
	Number       string    `json:"order"`
	Status       string    `json:"status"`
	Accrual      float64   `json:"accrual"`
	Processed_at time.Time `json:"processed_at"`
}
