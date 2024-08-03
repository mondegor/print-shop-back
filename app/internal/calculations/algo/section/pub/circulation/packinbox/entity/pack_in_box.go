package entity

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/packinbox"
)

const (
	ModelNameCirculationPackInBox = "public-api.Calculations.Algo.Circulation.PackInBox" // ModelNameCirculationPackInBox - название сущности
)

type (
	// ParsedData - разобранные валидные данные.
	ParsedData struct {
		Product packinbox.Product
		Box     packinbox.Box
	}

	// AlgoResult - результат работы алгоритма Algo.
	AlgoResult struct {
		FullBox          BoxResult `json:"fullBox"`
		RestBox          BoxResult `json:"restBox"`
		BoxesQuantity    uint64    `json:"boxesQuantity"`
		BoxesWeight      float64   `json:"boxesWeight"`
		ProductsVolume   float64   `json:"productsVolume"`
		BoxesVolume      float64   `json:"boxesVolume"`
		BoxesInnerVolume float64   `json:"boxesInnerVolume"`
	}

	// BoxResult - результаты вычислений параметров коробки.
	BoxResult struct {
		Weight              float64 `json:"weight"`
		Volume              float64 `json:"volume"`
		InnerVolume         float64 `json:"innerVolume"`
		ProductQuantity     uint64  `json:"productQuantity"`
		ProductVolume       float64 `json:"productVolume"`
		UnusedVolumePercent float64 `json:"unusedVolumePercent"`
	}
)
