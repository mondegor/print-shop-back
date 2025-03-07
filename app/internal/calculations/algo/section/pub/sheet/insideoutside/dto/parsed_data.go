package dto

import "github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2/rect2d"

type (
	// ParsedData - разобранные валидные данные.
	ParsedData struct {
		In  rect2d.Format
		Out rect2d.Format
	}
)
