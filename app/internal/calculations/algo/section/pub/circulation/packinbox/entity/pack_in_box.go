package entity

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/parallelepiped"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
)

const (
	ModelNameCirculationPackInBox = "public-api.Calculations.Algo.Circulation.PackInBox" // ModelNameCirculationPackInBox - название сущности
)

type (
	// RawData - сырые данные поступившие с обработчика.
	RawData struct {
		Product RawProduct
		Box     RawBox
	}

	// RawProduct - сырые данные об изделии.
	RawProduct struct {
		Format    string
		Thickness uint64
		WeightM2  uint64
		Quantity  uint64
	}

	// RawBox - сырые данные о коробке.
	RawBox struct {
		Format    string
		Thickness uint64
		Margins   string
		Weight    uint64
		MaxWeight uint64
	}

	// ParsedData - разобранные валидные данные.
	ParsedData struct {
		Product ParsedProduct
		Box     ParsedBox
	}

	// ParsedProduct - разобранные валидные данные об изделии.
	ParsedProduct struct {
		Format    rect.Format
		Thickness float64
		WeightM2  float64
		Quantity  uint64
	}

	// ParsedBox - разобранные валидные данные о коробке.
	ParsedBox struct {
		Format    parallelepiped.Format
		Thickness float64
		Margins   parallelepiped.Format
		Weight    float64
		MaxWeight float64
	}

	// AlgoResult - результат работы алгоритма Algo.
	AlgoResult struct {
		Box            BoxResult `json:"box"`
		LastBox        BoxResult `json:"lastBox"`
		ProductsVolume float64   `json:"productsVolume"`
		BoxesVolume    float64   `json:"boxesVolume"`
		BoxesQuantity  uint64    `json:"boxesQuantity"`
	}

	// BoxResult - результаты вычислений параметров коробки.
	BoxResult struct {
		ProductQuantity     uint64  `json:"productQuantity"`
		ProductVolume       float64 `json:"productVolume"`
		Weight              float64 `json:"weight"`
		Volume              float64 `json:"volume"`
		UnusedVolumePercent float64 `json:"unusedVolumePercent"`
	}
)
