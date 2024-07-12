package total

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
)

type (
	// AlgoTotal - вспомогательный алгоритм расчёта поверхности размещения элементов.
	AlgoTotal struct{}
)

// New - создаёт объект AlgoTotal.
func New() *AlgoTotal {
	return &AlgoTotal{}
}

// Calc - расчёт алгоритма.
func (ri *AlgoTotal) Calc(item rect.Item, out rect.Format, layouts base.Fragments) (layout rect.Format, restArea float64) {
	// рассчитывается основной блок элементов
	layout = ri.calcLayoutFormat(item, layouts[0])
	restArea = out.Area() - layout.Area()

	if len(layouts) > 1 {
		// обработка остаточного блока перевернутого на 90 градусов
		layout90Format := ri.calcLayoutFormat(item, layouts[1])
		restArea -= layout90Format.Area()

		blockTest := rect.Format{
			Width:  layout.Width + layout90Format.Height + item.Distance.Max(),
			Height: max(layout.Height, layout90Format.Width),
		}

		if blockTest.Width <= out.Width && blockTest.Height <= out.Height {
			layout = blockTest
		} else {
			layout = rect.Format{
				Width:  max(layout.Width, layout90Format.Height),
				Height: layout.Height + layout90Format.Width + item.Distance.Max(),
			}
		}

		restArea -= layout.Height * item.Distance.Max()
	}

	return layout, restArea
}

func (ri *AlgoTotal) calcLayoutFormat(item rect.Item, layout base.Fragment) rect.Format {
	format := rect.Format{
		Width:  float64(layout.ByWidth) * (item.Width + item.Distance.Width),
		Height: float64(layout.ByHeight) * (item.Height + item.Distance.Height),
	}

	// удаляется лишняя рамка элемента, т.к. она должна быть только между элементами
	return format.Diff(item.Distance)
}
