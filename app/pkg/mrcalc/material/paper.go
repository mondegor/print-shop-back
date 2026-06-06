package material

import (
	"print-shop-back/pkg/mrcalc/measure"
)

type (
	// Paper - материал "Бумага".
	Paper struct {
		Material
	}
)

// NewPaper - создаёт объект Paper.
func NewPaper(thickness measure.Micrometer, weightM2 measure.GramPerMeter2) *Paper {
	return &Paper{
		Material: Material{
			weightM2:  weightM2,
			thickness: thickness,
		},
	}
}
