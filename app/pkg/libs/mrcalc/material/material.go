package material

import (
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
)

type (
	// Material - абстрактный материал обладающий весом и толщеной.
	Material struct {
		weightM2  measure.GramPerMeter2
		thickness measure.Micrometer
	}
)

// NewMaterial - создаёт объект Material.
func NewMaterial(thickness measure.Micrometer, weightM2 measure.GramPerMeter2) *Material {
	return &Material{
		weightM2:  weightM2,
		thickness: thickness,
	}
}

// Weight - возвращает вес материала.
func (m *Material) Weight(quantity uint64, width, height measure.Micrometer) measure.GramPerMeter2 {
	return measure.GramPerMeter2(((quantity * uint64(m.weightM2) * uint64(width) / 1000) * uint64(height)) / 1000)
}

// Thickness - возвращает толщину материала.
func (m *Material) Thickness(quantity uint64) measure.Micrometer {
	return measure.Micrometer(quantity) * m.thickness
}
