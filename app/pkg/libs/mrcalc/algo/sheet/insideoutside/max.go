package insideoutside

import (
	"fmt"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/enum"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2/rect2d"
)

// AlgoMax - возвращает схему с максимальным количеством единиц указанного
// внутреннего формата, которое можно разместить во внешнем указанном формате.
func AlgoMax(in, out rect2d.Format) (rect2d.Fragments, error) {
	if !in.IsValid() {
		return nil, fmt.Errorf("in format is not valid: %s", in)
	}

	if !out.IsValid() {
		return nil, fmt.Errorf("out format is not valid: %s", out)
	}

	fragments := getMaxInsideOnOutside(in, out)
	if len(fragments) == 0 {
		return nil, nil
	}

	totalFirst := fragments[0].Layout.Quantity()
	total := fragments.TotalQuantity()

	// если формат не квадратный, то рассматривается
	// и книжный и альбомный варианты, а из них выбирается лучший
	if in.Width != in.Height {
		fragments90 := getMaxInsideOnOutside(in.Rotate90(), out)
		totalFirst90 := fragments90[0].Layout.Quantity()
		total90 := fragments90.TotalQuantity()

		// если кол-во элементов равно, то выбирается вариант, у которого
		// в первом фрагменте элементов больше ([{13, 3}] лучше, чем [{4, 9}, {3, 1}])
		// в случае, если у обоих вариантов элементов одинаковое кол-во в одном фрагменте,
		// то отдаётся предпочтение тому варианту, у которого элементов
		// на одной стороне больше ([{8, 3}] лучше, чем [{4, 6}])
		if total90 > total ||
			total90 == total && (totalFirst90 > totalFirst ||
				totalFirst90 == totalFirst && fragments90[0].Layout.Max() > fragments[0].Layout.Max()) {
			fragments = fragments90
		}
	}

	return fragments, nil
}

func getMaxInsideOnOutside(in, out rect2d.Format) rect2d.Fragments {
	layout := getQuantityInsideOnOutside(in, out)
	if layout.Quantity() == 0 {
		return nil
	}

	const maxFragments = 2
	fragments := make(rect2d.Fragments, 1, maxFragments)
	fragments[0] = rect2d.Fragment{
		Element:  in,
		Position: enum.PositionTop,
		Layout:   layout,
	}

	var (
		// остаток, на который ложится элемент повёрнутый на 90 градусов
		remaining rect2d.Format
		position  enum.Position
	)

	if in.Width > in.Height {
		position = enum.PositionOnside
		remaining.Width = out.Height
		remaining.Height = out.Width - in.Width*float64(layout.ByWidth)
	} else {
		position = enum.PositionBottom
		remaining.Width = out.Height - in.Height*float64(layout.ByHeight)
		remaining.Height = out.Width
	}

	if remaining.Width >= in.Width && remaining.Height >= in.Height {
		fragments = append(
			fragments,
			rect2d.Fragment{
				Element:  in.Rotate90(),
				Position: position,
				Layout:   getQuantityInsideOnOutside(in, remaining),
			},
		)
	}

	return fragments
}
