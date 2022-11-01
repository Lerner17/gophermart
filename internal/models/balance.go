package models

import (
	"database/sql"
	"encoding/json"
)

type Balance struct {
	Current   sql.NullFloat64 `json:"current"`
	Withdrawn sql.NullFloat64 `json:"withdrawn"`
}

func (b Balance) MarshalJSON() ([]byte, error) {
	var current float64
	var withdrawn int

	if !b.Current.Valid {
		current = 0
	} else {
		current = b.Current.Float64
	}

	if !b.Current.Valid {
		withdrawn = 0
	} else {
		withdrawn = int(b.Withdrawn.Float64)
	}
	return json.Marshal(struct {
		Current   float64 `json:"current"`
		Withdrawn int     `json:"withdrawn"`
	}{current, withdrawn})
}
