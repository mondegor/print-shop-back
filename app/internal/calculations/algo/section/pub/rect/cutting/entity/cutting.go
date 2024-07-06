package entity

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
)

const (
	ModelNameRectCutting = "public-api.Calculations.Algo.Rect.Cutting" // ModelNameRectCutting - название сущности
)

type (
	// RawData - сырые данные поступившие с обработчика.
	RawData struct {
		Fragments      []string
		DistanceFormat string
	}

	// ParsedData - разобранные валидные данные.
	ParsedData struct {
		Fragments      []base.Fragment
		DistanceFormat rect.Format
	}

	// QuantityResult - результат работы алгоритма AlgoQuantity.
	QuantityResult struct {
		Quantity int32 `json:"quantity"`
	}
)
