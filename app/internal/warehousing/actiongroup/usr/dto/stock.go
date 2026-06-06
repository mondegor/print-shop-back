package dto

import (
	"github.com/google/uuid"

	"print-shop-back/internal/warehousing/xtype"
)

type (
	// AddMoreStockContainer - comment struct.
	AddMoreStockContainer struct {
		AccountID   uuid.UUID
		ContainerID uint64
		LocationID  uint64
		Quantity    int
	}

	// MoveStockContainer - comment struct.
	MoveStockContainer struct {
		AccountID  uuid.UUID
		StockID    uint64
		LocationID uint64
		Quantity   int
	}

	// TransferStockContainers - comment struct.
	TransferStockContainers struct {
		AccountID uuid.UUID
		Stocks    []StockQuantity
	}

	// StockQuantity - comment struct.
	StockQuantity struct {
		StockID  uint64
		Quantity int
	}

	// StockParams - comment struct.
	StockParams struct {
		AccountID uuid.UUID
		Filter    StockListFilter
		Cursor    xtype.StockCursor
	}

	// StockListFilter - comment struct.
	StockListFilter struct {
		SearchLocations  []uint64
		SearchContainers []uint64
	}
)
