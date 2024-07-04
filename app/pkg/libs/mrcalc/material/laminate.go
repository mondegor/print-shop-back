package material

import (
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
)

type (
	Laminate struct {
		Material
	}
)

func NewLaminate(thickness measure.Micrometer, weightM2 measure.GramsPerMeter2) *Laminate {
	return &Laminate{
		Material: Material{
			weightM2: weightM2,
			thickness: thickness,
		},
	}
}
