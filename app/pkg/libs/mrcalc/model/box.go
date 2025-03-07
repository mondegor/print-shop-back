package model

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s3/rect3d"
)

const (
	twoBoxSides = 2 // стенки коробки с двух сторон
)

type (
	// Box - модель коробки для размещения в ней изделий.
	Box struct {
		Format    rect3d.Format
		Thickness float64
		Margins   rect3d.Format
		Weight    float64
		MaxWeight float64
	}
)

// InnerFormat - возвращает внутренний формат коробки.
func (m Box) InnerFormat() rect3d.Format {
	// вычитается толщина стенок
	return m.Format.SubSame(m.Thickness * twoBoxSides)
}

// NarrowedFormat - возвращает внутренний формат коробки с учётом рамок.
func (m Box) NarrowedFormat() rect3d.Format {
	boxSidesThickness := m.Thickness * twoBoxSides

	// вычитается толщина стенок и указанная рамка
	return m.Format.Sub(
		rect3d.Format{
			Length: boxSidesThickness + m.Margins.Length,
			Width:  boxSidesThickness + m.Margins.Width,
			Height: boxSidesThickness + m.Margins.Height,
		},
	)
}

// Volume - возвращает объём коробки.
func (m Box) Volume() float64 {
	return m.Format.Volume()
}

// InnerVolume - возвращает внутренний объём коробки.
func (m Box) InnerVolume() float64 {
	return m.InnerFormat().Volume()
}
