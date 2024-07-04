package material

import (
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
)

type (
	Paper struct {
		Material
	}
)

func NewPaper(thickness measure.Micrometer, weightM2 measure.GramsPerMeter2) *Paper {
	return &Paper{
		Material: Material{
			weightM2: weightM2,
			thickness: thickness,
		},
	}
}
