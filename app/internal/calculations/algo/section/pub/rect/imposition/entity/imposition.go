package entity

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/imposition"
)

const (
	ModelNameRectImposition = "public-api.Calculations.Rect.Imposition" // ModelNameRectImposition - название сущности
)

type (
	// RawData - сырые данные поступившие с обработчика,
	// которые предназначены для вычисления алгоритма.
	RawData struct {
		ItemFormat       string
		ItemBorderFormat string
		OutFormat        string
		AllowRotation    bool
		UseMirror        bool
	}

	// ParsedData - разобранные валидные данные поступившие с обработчика,
	// которые предназначены для вычисления алгоритма.
	ParsedData struct {
		Item rect.Item
		Out  rect.Format
		Opts imposition.Options
	}
)
