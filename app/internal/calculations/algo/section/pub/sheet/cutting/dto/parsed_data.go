package dto

import (
	"print-shop-back/pkg/mrcalc/s2/rect2d"
)

type (
	// ParsedData - разобранные валидные данные.
	ParsedData struct {
		Fragments      []rect2d.Layout
		DistanceFormat rect2d.Format
	}
)
