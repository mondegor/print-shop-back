package insideoutside

import (
	"fmt"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
)

// AlgoQuantity - возвращает количество единиц указанного
// внутреннего формата, которое можно разместить по вертикали и горизонтали
// во внешнем указанном формате (без использования поворотов).
func AlgoQuantity(in, out rect.Format) (base.Fragment, error) {
	if in.Width <= 0 {
		return base.Fragment{}, fmt.Errorf("in.Width is zero or negative: %.2f", in.Width)
	}

	if in.Height <= 0 {
		return base.Fragment{}, fmt.Errorf("in.Height is zero or negative: %.2f", in.Height)
	}

	return getQuantityInsideOnOutside(in, out), nil
}

// getQuantityInsideOnOutside - необходимо гарантировать in.Width > 0 и in.Height > 0.
func getQuantityInsideOnOutside(in, out rect.Format) base.Fragment {
	return base.Fragment{
		ByWidth:  uint64(mrlib.RoundFloat8(out.Width / in.Width)),
		ByHeight: uint64(mrlib.RoundFloat8(out.Height / in.Height)),
	}
}
