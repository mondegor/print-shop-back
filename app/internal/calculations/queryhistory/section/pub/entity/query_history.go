package entity

import "time"

const (
	// ModelNameQueryHistory - название сущности.
	ModelNameQueryHistory = "public-api.Calculations.QueryHistory"
)

type (
	// QueryHistoryItem - comment struct.
	QueryHistoryItem struct {
		Caption   string    `json:"caption"`
		Params    string    `json:"params"`
		Result    string    `json:"result"`
		CreatedAt time.Time `json:"createdAt"`
	}
)
