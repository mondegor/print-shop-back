package entity

import (
	"time"

	"print-shop-back/internal/warehousing/enum/activitystatus"
	"print-shop-back/internal/warehousing/xtype"
)

type (
	// Territory - comment struct.
	Territory struct { // DB: warehousing.stores
		ID          uint64              `json:"id"`
		TagVersion  uint32              `json:"version"`
		AccountID   uint64              `json:"-"`
		Caption     string              `json:"caption"`
		Address     xtype.Volume        `json:"address"`
		Description activitystatus.Enum `json:"description"`
		CodePattern string              `json:"code_pattern"`
		CreatedAt   time.Time           `json:"created_at"`
		UpdatedAt   time.Time           `json:"updated_at"`
	}

	// TerritoryState - comment struct.
	TerritoryState struct { // DB: warehousing.stores
		ID          uint64
		TagVersion  uint32
		CodePattern string
	}
)
