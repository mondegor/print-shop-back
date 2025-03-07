package rect2d

import "github.com/mondegor/print-shop-back/pkg/libs/mrcalc/enum"

type (
	// Fragments - comment struct.
	Fragments []Fragment
)

// TotalQuantity - возвращает общее кол-во элементов размещённых во фрагментах.
func (f *Fragments) TotalQuantity() (total uint64) {
	for _, fragment := range *f {
		total += fragment.Layout.Quantity()
	}

	return total
}

// FragmentDistance - возвращает дистанцию между фрагментами.
func (f *Fragments) FragmentDistance() (max float64) {
	if len(*f) > 1 {
		for _, fragment := range *f {
			if fragment.Position == enum.PositionTop {
				return fragment.Distance.Max()
			}
		}
	}

	return 0
}

// Round - возвращает округлённые фрагменты.
func (f *Fragments) Round() Fragments {
	fragments := make([]Fragment, 0, len(*f))

	for _, fragment := range *f {
		fragments = append(fragments, fragment.Round())
	}

	return fragments
}
