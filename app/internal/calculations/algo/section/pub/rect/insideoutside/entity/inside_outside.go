package entity

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
)

const (
	ModelNameInsideOutside = "public-api.Calculations.Rect.InsideOutside" // ModelNameInsideOutside - название сущности
)

type (
	// RawData - сырые данные поступившие с обработчика,
	// которые предназначены для вычисления алгоритма.
	RawData struct {
		InFormat  string
		OutFormat string
	}

	// ParsedData - разобранные валидные данные поступившие с обработчика,
	// которые предназначены для вычисления алгоритма.
	ParsedData struct {
		In  rect.Format
		Out rect.Format
	}

	// AlgoQuantityResult - разобранные валидные данные поступившие с обработчика,
	// которые предназначены для вычисления алгоритма.
	AlgoQuantityResult struct {
		Fragment base.Fragment `json:"fragment"`
		Total    uint64        `json:"total"`
	}

	// AlgoMaxResult - разобранные валидные данные поступившие с обработчика,
	// которые предназначены для вычисления алгоритма.
	AlgoMaxResult struct {
		Fragments []base.Fragment `json:"fragments"`
		Total     uint64          `json:"total"`
	}
)
