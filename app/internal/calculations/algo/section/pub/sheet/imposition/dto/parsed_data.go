package dto

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/algo/sheet/imposition"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2/rect2d"
)

type (
	// ParsedData - разобранные валидные данные.
	ParsedData struct {
		Element  rect2d.Format
		Distance rect2d.Format
		Out      rect2d.Format
		Opts     imposition.Options
	}
)
