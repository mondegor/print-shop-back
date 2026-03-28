package httpv1

import (
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/entity"
)

type (
	// AddMoreStockRequest - comment struct.
	AddMoreStockRequest struct {
		ContainerID uint64 `json:"container_id"`
		LocationID  uint64 `json:"location_id"`
		Quantity    int    `json:"quantity"`
	}

	// MoveStockRequest - comment struct.
	MoveStockRequest struct {
		StockID    uint64 `json:"stock_id"`
		LocationID uint64 `json:"location_id"`
		Quantity   int    `json:"quantity,omitempty"`
	}

	// TransferStockRequest - comment struct.
	TransferStockRequest struct {
		StockID  uint64 `json:"stock_id"`
		Quantity int    `json:"quantity"`
	}

	// StockListResponse - comment struct.
	StockListResponse struct {
		Items   []entity.Stock `json:"stocks"`
		Cursor  string         `json:"cursor"`
		HasNext bool           `json:"has_next"`
	}

	// SuccessAddedStockResponse - comment struct.
	SuccessAddedStockResponse struct {
		StockID uint64 `json:"stock_id"`
	}

	// SuccessMovedStockResponse - comment struct.
	SuccessMovedStockResponse struct {
		StockID uint64 `json:"stock_id"`
	}
)
