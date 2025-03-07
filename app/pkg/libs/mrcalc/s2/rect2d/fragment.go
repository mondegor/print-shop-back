package rect2d

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/enum"
)

type (
	// Fragment - фрагмент прямоугольного формата на котором размещены
	// элементы указанного формата на указанной дистанции между ними.
	Fragment struct {
		Element  Format        `json:"element"`
		Distance Format        `json:"distance"` // расстояние между элементами по вертикали и горизонтали
		Layout   Layout        `json:"layout"`
		Position enum.Position `json:"position"`
	}
)

// ElementWithDistance - возвращает формат элемента вместе с расстоянием
// до соседних элементов с двух перпендикулярных сторон.
func (f *Fragment) ElementWithDistance() Format {
	return f.Element.Add(f.Distance)
}

// Format - возвращает формат фрагмента.
func (f *Fragment) Format() Format {
	return Format{
		Width:  f.Width(),
		Height: f.Height(),
	}
}

// Width - возвращает ширину фрагмента.
func (f *Fragment) Width() float64 {
	if f.Layout.ByWidth == 0 {
		return 0
	}

	return (f.Element.Width+f.Distance.Width)*float64(f.Layout.ByWidth) - f.Distance.Width
}

// Height - возвращает высоту фрагмента.
func (f *Fragment) Height() float64 {
	if f.Layout.ByHeight == 0 {
		return 0
	}

	return (f.Element.Height+f.Distance.Height)*float64(f.Layout.ByHeight) - f.Distance.Height
}

// Round - возвращает округлённый фрагмент.
func (f *Fragment) Round() Fragment {
	return Fragment{
		Element:  f.Element.Round(),
		Distance: f.Distance.Round(),
		Layout:   f.Layout,
		Position: f.Position,
	}
}
