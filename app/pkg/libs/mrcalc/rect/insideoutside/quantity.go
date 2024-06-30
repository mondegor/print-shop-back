package insideoutside

import (
	"fmt"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
)

// AlgoQuantity - возвращает количество единиц указанного
// внутреннего формата, которое можно разместить по вертикали и горизонтали
// во внешнем указанном формате (без использования поворотов).
func AlgoQuantity(in, out rect.Format) (base.Fragment, error) {
	if in.Width < 1 {
		return base.Fragment{}, fmt.Errorf("in.Width is zero or negative: %d", in.Width)
	}

	if in.Height < 1 {
		return base.Fragment{}, fmt.Errorf("in.Height is zero or negative: %d", in.Height)
	}

	return getQuantityInsideOnOutside(in, out), nil
}

// getQuantityInsideOnOutside - необходимо гарантировать in.Width > 0 и in.Height > 0.
func getQuantityInsideOnOutside(in, out rect.Format) base.Fragment {
	return base.Fragment{
		ByWidth:  uint64(out.Width / in.Width),
		ByHeight: uint64(out.Height / in.Height),
	}
}
