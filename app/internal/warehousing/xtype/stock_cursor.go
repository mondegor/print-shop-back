package xtype

import (
	"strconv"

	"github.com/mondegor/go-sysmess/mrtype"
)

type (
	// StockCursor - параметры для выборки части списка элементов.
	StockCursor struct {
		StockID uint64
		Limit   int
	}
)

// NewStockCursor - создаёт объект StockCursor.
func NewStockCursor(params mrtype.CursorParams) StockCursor {
	if params.Limit == 0 {
		params.Limit = 1
	}

	id, err := strconv.ParseUint(params.Value, 10, 64)
	if err != nil {
		return StockCursor{
			Limit: params.Limit,
		}
	}

	return StockCursor{
		StockID: id,
		Limit:   params.Limit,
	}
}
