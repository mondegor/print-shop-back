package dto

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/model"
)

type (
	// ParsedData - разобранные валидные данные.
	ParsedData struct {
		SheetHeap       model.SheetStack
		QuantityInStack uint64
	}
)
