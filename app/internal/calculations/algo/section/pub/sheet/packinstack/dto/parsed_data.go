package dto

import (
	"print-shop-back/pkg/mrcalc/model"
)

type (
	// ParsedData - разобранные валидные данные.
	ParsedData struct {
		SheetHeap       model.SheetStack
		QuantityInStack uint64
	}
)
