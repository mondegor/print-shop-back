package entity

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
)

const (
	ModelNameRectInsideOutside = "public-api.Calculations.Algo.Rect.InsideOutside" // ModelNameRectInsideOutside - название сущности
)

type (
	// ParsedData - разобранные валидные данные.
	ParsedData struct {
		In  rect.Format
		Out rect.Format
	}

	// AlgoQuantityResult - результат работы алгоритма AlgoQuantity.
	AlgoQuantityResult struct {
		Fragment base.Fragment `json:"fragment"`
		Total    uint64        `json:"total"`
	}

	// AlgoMaxResult - результат работы алгоритма AlgoMax.
	AlgoMaxResult struct {
		Fragments []base.Fragment `json:"fragments"`
		Total     uint64          `json:"total"`
	}
)
