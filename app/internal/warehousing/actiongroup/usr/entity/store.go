package entity

import (
	"time"

	"github.com/google/uuid"

	"github.com/mondegor/print-shop-back/internal/warehousing/enum/activitystatus"
	"github.com/mondegor/print-shop-back/internal/warehousing/enum/storekind"
	"github.com/mondegor/print-shop-back/internal/warehousing/xtype"
)

type (
	// Store - comment struct.
	Store struct { // DB: warehousing.stores
		ID               uint64              `json:"id"`
		TagVersion       uint32              `json:"version"`
		AccountID        uuid.UUID           `json:"-"`
		TerritoryID      uint64              `json:"territory_id"`
		Kind             storekind.Enum      `json:"kind"`
		Code             string              `json:"code"`
		Volume           xtype.Volume        `json:"volume"`
		Status           activitystatus.Enum `json:"status"`
		ContainersVolume float64             `json:"containers_volume"`
		CreatedAt        time.Time           `json:"created_at"`
		UpdatedAt        time.Time           `json:"updated_at"`
		Fullness         float64             `json:"fullness"` // TODO: временно, перенести
	}

	// StoreState - comment struct.
	StoreState struct { // DB: warehousing.stores
		ID          uint64
		TagVersion  uint32
		TerritoryID uint64
		Kind        storekind.Enum
		Code        string
		Status      activitystatus.Enum
	}
)

// CreateStoreCursorValue - comment func.
func CreateStoreCursorValue(items []Store) string {
	if len(items) == 0 {
		return ""
	}

	return items[len(items)-1].Code
}
