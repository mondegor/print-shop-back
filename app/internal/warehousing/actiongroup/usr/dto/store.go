package dto

import (
	"github.com/google/uuid"

	"print-shop-back/internal/warehousing/xtype"
)

type (
	// StoreParams - comment struct.
	StoreParams struct {
		AccountID uuid.UUID
		Filter    StoreListFilter
		Cursor    xtype.StoreCursor
	}

	// StoreListFilter - comment struct.
	StoreListFilter struct {
		SearchCode        string
		SearchTerritories []uint64
	}
)
