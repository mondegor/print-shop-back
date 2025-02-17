package entity

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/imposition"
)

const (
	ModelNameRectImposition = "public-api.Calculations.Algo.Rect.Imposition" // ModelNameRectImposition - название сущности
)

type (
	// ParsedData - разобранные валидные данные.
	ParsedData struct {
		Item rect.Item
		Out  rect.Format
		Opts imposition.Options
	}

	// AlgoResult - результат вычислений спуска полос.
	AlgoResult struct {
		Layout    rect.Format     `json:"layout"`
		Item      rect.Format     `json:"item"`
		Distance  rect.Format     `json:"distance"`
		Fragments []base.Fragment `json:"fragments"`
		Total     uint64          `json:"total"`
		Garbage   float64         `json:"garbage"`
	}

	// AlgoVariantResult - .
	AlgoVariantResult struct {
		Name   string     `json:"name"`
		Result AlgoResult `json:"result"`
	}
)
