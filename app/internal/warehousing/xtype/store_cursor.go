package xtype

import (
	"github.com/mondegor/go-core/mrtype"
)

type (
	// StoreCursor - параметры для выборки части списка элементов.
	StoreCursor struct {
		Code  string
		Limit int
	}
)

// NewStoreCursor - создаёт объект StoreCursor.
func NewStoreCursor(params mrtype.CursorParams) StoreCursor {
	if params.Limit == 0 {
		params.Limit = 1
	}

	return StoreCursor{
		Code:  params.Value,
		Limit: params.Limit,
	}
}
