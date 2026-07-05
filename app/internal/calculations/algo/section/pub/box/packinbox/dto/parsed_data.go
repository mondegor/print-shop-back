package dto

import (
	"print-shop-back/pkg/mrcalc/model"
)

type (
	// ParsedData - разобранные валидные данные.
	ParsedData struct {
		ProductHeap model.ProductStack
		Box         model.Box
	}
)
