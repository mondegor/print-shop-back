package entity

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
)

const (
	ModelNameRectCutting = "public-api.Calculations.Algo.Rect.Cutting" // ModelNameRectCutting - название сущности
)

type (
	// ParsedData - разобранные валидные данные.
	ParsedData struct {
		Fragments      []base.Fragment
		DistanceFormat rect.Format
	}

	// QuantityResult - результат работы алгоритма AlgoQuantity.
	QuantityResult struct {
		Quantity uint64 `json:"quantity"`
	}
)
