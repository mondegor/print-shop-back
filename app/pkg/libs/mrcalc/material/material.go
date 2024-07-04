package material

import (
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
)

type (
	Material struct {
		weightM2 measure.GramsPerMeter2
		thickness measure.Micrometer
	}
)

func NewMaterial(thickness measure.Micrometer, weightM2 measure.GramsPerMeter2) *Material {
	return &Material{
		weightM2: weightM2,
		thickness: thickness,
	}
}

func (m *Material) Thickness(quantity uint64) measure.Micrometer {
	return measure.Micrometer(quantity) * m.thickness
}

func (m *Material) Weight(quantity uint64, width, height measure.Micrometer) measure.GramsPerMeter2 {
	return measure.GramsPerMeter2(((quantity * uint64(m.weightM2) * uint64(width) / 1000) * uint64(height)) / 1000)
}
