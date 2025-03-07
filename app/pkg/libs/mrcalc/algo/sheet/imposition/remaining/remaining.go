package remaining

import (
	"fmt"

	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/algo/sheet/insideoutside"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/enum"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2/rect2d"
)

type (
	// AlgoRemaining - вспомогательный алгоритм расчёта остатка.
	AlgoRemaining struct {
		logger mrlog.Logger
	}
)

// New - создаёт объект AlgoRemaining.
func New(logger mrlog.Logger) *AlgoRemaining {
	return &AlgoRemaining{
		logger: logger,
	}
}

// Calc - расчёт алгоритма.
func (ri *AlgoRemaining) Calc(inFragment rect2d.Fragment, out rect2d.Format) (rf rect2d.Fragment, err error) {
	out, position := ri.remainingFormat(inFragment, out)

	if !out.IsValid() {
		return rect2d.Fragment{}, nil
	}

	item90WithDistance := inFragment.Element.Rotate90().Add(inFragment.Distance.Rotate90())

	layout, err := insideoutside.AlgoQuantity(item90WithDistance, out)
	if err != nil {
		return rect2d.Fragment{}, fmt.Errorf("AlgoRemaining.Calc[in=%+v, outRemaining=%s]: %w", item90WithDistance, out, err)
	}

	// если ни один элемент невозможно разместить в остаточном формате
	if layout.Quantity() == 0 {
		ri.logger.Debug().MsgFunc(
			func() string {
				return fmt.Sprintf(
					"- skipped: item format %s not on remaining format %s with fict margins",
					item90WithDistance,
					out,
				)
			},
		)

		return rect2d.Fragment{}, nil
	}

	ri.logger.Debug().MsgFunc(
		func() string {
			return fmt.Sprintf(
				"- placed item %s on remaining out format %s with fict margins: %s, %s, %d * %d = %d",
				item90WithDistance,
				out,
				item90WithDistance.OrientationType(),
				position,
				inFragment.Layout.ByWidth,
				inFragment.Layout.ByHeight,
				inFragment.Layout.Quantity(),
			)
		},
	)

	return rect2d.Fragment{
		Element:  inFragment.Element.Rotate90(),
		Distance: inFragment.Distance.Rotate90(),
		Layout:   layout,
		Position: position,
	}, nil
}

func (ri *AlgoRemaining) remainingFormat(inFragment rect2d.Fragment, out rect2d.Format) (format rect2d.Format, position enum.Position) {
	element := inFragment.ElementWithDistance()

	if element.Width >= element.Height {
		position = enum.PositionOnside

		// расстояние по ширине между вертикально ориентированными элементами и горизонтально ориентированными
		inFragmentWidth := element.Width*float64(inFragment.Layout.ByWidth) + inFragment.Distance.Max() - inFragment.Distance.Width

		format = rect2d.Format{
			Width:  out.Width - inFragmentWidth,
			Height: out.Height,
		}
	} else {
		position = enum.PositionBottom

		// расстояние по высоте между вертикально ориентированными элементами и горизонтально ориентированными
		inFragmentHeight := element.Height*float64(inFragment.Layout.ByHeight) + inFragment.Distance.Max() - inFragment.Distance.Height

		format = rect2d.Format{
			Width:  out.Width,
			Height: out.Height - inFragmentHeight,
		}
	}

	// прибавляется фиктивная граница для учёта размещения граничных элементов
	return format.Add(inFragment.Distance.Rotate90()), position
}
