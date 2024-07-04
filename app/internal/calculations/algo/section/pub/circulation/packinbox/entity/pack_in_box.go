package entity

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
)

const (
	ModelNamePackInBox = "public-api.Calculations.Circulation.PackInBox" // ModelNamePackInBox - название сущности
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
