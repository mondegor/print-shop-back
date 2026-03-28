package entity

import (
	"strconv"
	"time"

	"github.com/google/uuid"
)

type (
	// Stock - comment struct.
	Stock struct { // DB: warehousing.container_stocks
		ID                uint64    `json:"id"`
		AccountID         uuid.UUID `json:"-"`
		ContainerID       uint64    `json:"container_id"`
		LocationID        uint64    `json:"location_id"`
		ContainerQuantity int       `json:"container_quantity"`
		ContainerVolume   float64   `json:"container_volume"` // m3
		CreatedAt         time.Time `json:"created_at"`
	}
)

// CreateStockCursorValue - comment func.
func CreateStockCursorValue(items []Stock) string {
	if len(items) == 0 {
		return ""
	}

	return strconv.FormatUint(items[len(items)-1].ID, 10)
}
