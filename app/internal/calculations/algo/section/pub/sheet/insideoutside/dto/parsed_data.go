package dto

import "print-shop-back/pkg/mrcalc/s2/rect2d"

type (
	// ParsedData - разобранные валидные данные.
	ParsedData struct {
		In  rect2d.Format
		Out rect2d.Format
	}
)
