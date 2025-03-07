package total

import "github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2/rect2d"

type (
	// AlgoTotal - вспомогательный алгоритм расчёта поверхности размещения элементов.
	AlgoTotal struct{}
)

// New - создаёт объект AlgoTotal.
func New() *AlgoTotal {
	return &AlgoTotal{}
}

// Calc - расчёт алгоритма.
func (ri *AlgoTotal) Calc(inFragments rect2d.Fragments, out rect2d.Format) (containerFormat rect2d.Format, restArea float64) {
	if len(inFragments) == 0 {
		return rect2d.Format{}, 0
	}

	// рассчитывается основной блок элементов
	containerFormat = inFragments[0].Format()
	restArea = out.Area() - containerFormat.Area()

	if len(inFragments) > 1 {
		// обработка остаточного блока
		remainingFormat := inFragments[1].Format()
		restArea -= remainingFormat.Area()
		elementDistanceMax := inFragments[0].Distance.Max()

		blockTest := rect2d.Format{
			Width:  containerFormat.Width + remainingFormat.Width + elementDistanceMax,
			Height: max(containerFormat.Height, remainingFormat.Height),
		}

		if blockTest.Width <= out.Width && blockTest.Height <= out.Height {
			containerFormat = blockTest
		} else {
			containerFormat = rect2d.Format{
				Width:  max(containerFormat.Width, remainingFormat.Width),
				Height: containerFormat.Height + remainingFormat.Height + elementDistanceMax,
			}
		}

		restArea -= containerFormat.Height * elementDistanceMax
	}

	return containerFormat, restArea
}
