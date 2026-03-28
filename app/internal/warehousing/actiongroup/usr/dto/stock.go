package dto

import (
	"github.com/google/uuid"

	"github.com/mondegor/print-shop-back/internal/warehousing/enum/locationkind"
	"github.com/mondegor/print-shop-back/internal/warehousing/xtype"
)

type (
	// CreateStockContainer - comment struct.
	CreateStockContainer struct {
		AccountID  uuid.UUID
		Kind       locationkind.Enum
		Code       string
		Volume     xtype.Volume
		Tags       []string
		Images     []string
		LocationID uint64
		Quantity   int
	}

	// CreateStockContainerResult - comment struct.
	CreateStockContainerResult struct {
		ID      uint64
		Code    string
		Marker  uint16
		StockID uint64
	}

	// MoveStockContainer - comment struct.
	MoveStockContainer struct {
		AccountID  uuid.UUID
		StockID    uint64
		LocationID uint64
		Quantity   int
	}

	// AddMoreStockContainer - comment struct.
	AddMoreStockContainer struct {
		AccountID   uuid.UUID
		ContainerID uint64
		LocationID  uint64
		Quantity    int
	}

	// ChangeStockLocation - comment struct.
	ChangeStockLocation struct {
		AccountID    uuid.UUID
		StockID      uint64
		LocationID   uint64
		ToLocationID uint64
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
