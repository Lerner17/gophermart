package models

import "time"

type Withdraw struct {
	Number       string    `json:"number" db:"order"`
	Processed_at time.Time `json:"processed_at"`
	Sum          float64   `json:"sum"`
}
