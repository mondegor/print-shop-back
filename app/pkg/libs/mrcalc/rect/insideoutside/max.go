package insideoutside

import (
	"fmt"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
)

// AlgoMax - возвращает максимальное количество единиц указанного
// внутреннего формата, которое можно разместить во внешнем указанном формате.
func AlgoMax(in, out rect.Format) (base.Fragments, error) {
	if !in.IsValid() {
		return base.Fragments{}, fmt.Errorf("in format is not valid: %s", in)
	}

	if !out.IsValid() {
		return base.Fragments{}, fmt.Errorf("out format is not valid: %s", out)
	}

	out = out.Cast()

	fragments := getMaxInsideOnOutside(in, out)
	if len(fragments) == 0 {
		return nil, nil
	}

	totalFirst := fragments[0].Total()
	total := fragments.Total()

	// если формат не квадратный то рассматривается
	// и книжный и альбомный варианты, а из них выбирается лучший
	if in.Width != in.Height {
		fragments90 := getMaxInsideOnOutside(in.Change(), out)
		totalFirst90 := fragments90[0].Total()
		total90 := fragments90.Total()

		// если кол-во элементов равно, то выбирается вариант,
		// у которого в первом фрагменте элементов больше ([{13, 3}] лучше, чем [{4, 9}, {3, 1}])
		// в случае, если у обоих вариантов элементов одинаковое кол-во в одном фрагменте,
		// то отдаётся предпочтение тому варианту, у которого элементов
		// на одной стороне больше ([{8, 3}] лучше, чем [{4, 6}])
		if total90 > total ||
			total90 == total && (totalFirst90 > totalFirst ||
				totalFirst90 == totalFirst && fragments90[0].Max() > fragments[0].Max()) {
			fragments = fragments90
		}
	}

	return fragments, nil
}

func getMaxInsideOnOutside(in, out rect.Format) base.Fragments {
	fragment := getQuantityInsideOnOutside(in, out)
	if fragment.Total() == 0 {
		return nil
	}

	const maxFragments = 2
	fragments := make(base.Fragments, 0, maxFragments)
	fragments = append(fragments, fragment)
	remaining := rect.Format{}

	if in.Width > in.Height {
		remaining.Width = out.Height
		remaining.Height = out.Width - in.Width*int64(fragment.ByWidth)
	} else {
		remaining.Width = out.Height - in.Height*int64(fragment.ByHeight)
		remaining.Height = out.Width
	}

	if remaining.Width >= in.Width && remaining.Height >= in.Height {
		fragments = append(fragments, getQuantityInsideOnOutside(in, remaining))
	}

	return fragments
}
