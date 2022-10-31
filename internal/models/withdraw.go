package models

import "time"

type Withdraw struct {
	Number      string    `json:"order"`
	ProcessedAt time.Time `json:"processed_at"`
	Sum         float64   `json:"sum"`
}
