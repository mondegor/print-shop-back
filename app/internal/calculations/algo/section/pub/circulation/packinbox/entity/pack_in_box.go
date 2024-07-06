package entity

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
)

const (
	ModelNameCirculationPackInBox = "public-api.Calculations.Algo.Circulation.PackInBox" // ModelNameCirculationPackInBox - название сущности
)

type (
	// RawData - сырые данные поступившие с обработчика.
	RawData struct {
		Format string
	}

	// ParsedData - разобранные валидные данные.
	ParsedData struct {
		Format rect.Format
	}

	// AlgoResult - результат работы алгоритма Algo.
	AlgoResult struct {
		Format rect.Format `json:"format"`
	}
)
