package total

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
)

type (
	// AlgoTotal - спуск полос.
	// * Рассчитывается максимальное количество элементов указанного формата,
	// * которое можно разместить на листе указанного формата.
	// * Также учитывается расстояние между элементами по горизонтали и вертикали.
	AlgoTotal struct{}
)

// New - создаёт объект AlgoTotal.
// Поддерживается параметр allowRotation при true разрешается
// располагать элементы повёрнутые на 90 градусов к друг другу.
func New() *AlgoTotal {
	return &AlgoTotal{}
}

// Calc - слайс layouts всегда содержит минимум 1 элемент и не больше 2‑х элементов.
func (ri *AlgoTotal) Calc(item rect.Item, out rect.Format, layouts base.Fragments) (layout rect.Format, restArea int64) {
	// рассчитывается основной блок элементов
	layout = ri.calcLayoutFormat(item, layouts[0])
	restArea = out.Area() - layout.Area()

	if len(layouts) > 1 {
		// обработка остаточного блока перевернутого на 90 градусов
		layout90Format := ri.calcLayoutFormat(item, layouts[1])
		restArea -= layout90Format.Area()

		blockTest := rect.Format{
			Width:  layout.Width + layout90Format.Height + item.Border.Max(),
			Height: max(layout.Height, layout90Format.Width),
		}

		if blockTest.Width <= out.Width && blockTest.Height <= out.Height {
			layout = blockTest
		} else {
			layout = rect.Format{
				Width:  max(layout.Width, layout90Format.Height),
				Height: layout.Height + layout90Format.Width + item.Border.Max(),
			}
		}

		restArea -= layout.Height * item.Border.Max()
	}

	return layout, restArea
}

func (ri *AlgoTotal) calcLayoutFormat(item rect.Item, layout base.Fragment) rect.Format {
	format := rect.Format{
		Width:  int64(layout.ByWidth) * (item.Width + item.Border.Width),
		Height: int64(layout.ByHeight) * (item.Height + item.Border.Height),
	}

	// удаляется лишняя рамка элемента, т.к. она должна быть только между элементами
	return format.Diff(item.Border)
}
