package material

import (
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
)

type (
	// Laminate - материал "Ламинат".
	Laminate struct {
		Material
	}
)

// NewLaminate - создаёт объект Laminate.
func NewLaminate(thickness measure.Micrometer, weightM2 measure.GramPerMeter2) *Laminate {
	return &Laminate{
		Material: Material{
			weightM2:  weightM2,
			thickness: thickness,
		},
	}
}
