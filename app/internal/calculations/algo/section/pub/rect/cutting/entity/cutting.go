package entity

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
)

const (
	ModelNameCutting = "public-api.Calculations.Rect.Cutting" // ModelNameCutting - название сущности
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

	// AlgoQuantityResult - результат работы алгоритма AlgoQuantity.
	AlgoQuantityResult struct {
		Quantity uint64 `json:"quantity"`
	}
)
