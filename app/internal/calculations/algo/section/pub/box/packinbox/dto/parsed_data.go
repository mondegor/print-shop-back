package dto

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/model"
)

type (
	// ParsedData - разобранные валидные данные.
	ParsedData struct {
		ProductHeap model.ProductStack
		Box         model.Box
	}
)
