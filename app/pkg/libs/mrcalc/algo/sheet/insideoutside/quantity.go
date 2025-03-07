package insideoutside

import (
	"fmt"

	"github.com/mondegor/go-webcore/mrlib"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2/rect2d"
)

// AlgoQuantity - возвращает схему с количеством единиц указанного
// внутреннего формата, которое можно разместить по вертикали и горизонтали
// во внешнем указанном формате (без использования поворотов).
func AlgoQuantity(in, out rect2d.Format) (rect2d.Layout, error) {
	if in.Width <= 0 {
		return rect2d.Layout{}, fmt.Errorf("in.Width is zero or negative: %.2f", in.Width)
	}

	if in.Height <= 0 {
		return rect2d.Layout{}, fmt.Errorf("in.Height is zero or negative: %.2f", in.Height)
	}

	return getQuantityInsideOnOutside(in, out), nil
}

// getQuantityInsideOnOutside - необходимо гарантировать in.Width > 0 и in.Height > 0.
func getQuantityInsideOnOutside(in, out rect2d.Format) rect2d.Layout {
	return rect2d.Layout{
		ByWidth:  uint64(mrlib.RoundFloat8(out.Width / in.Width)),
		ByHeight: uint64(mrlib.RoundFloat8(out.Height / in.Height)),
	}
}
