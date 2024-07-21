package remaining

import (
	"fmt"

	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/insideoutside"
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
func (ri *AlgoRemaining) Calc(layout base.Fragment, item rect.Item, out rect.Format) (base.Fragment, error) {
	outRemaining, remainingPosition := ri.remainingFormat(item, out, layout)

	if outRemaining.Area() == 0 {
		return base.Fragment{}, nil
	}

	remainingLayout, err := insideoutside.AlgoQuantity(item.WithDistance(), outRemaining)
	if err != nil {
		return base.Fragment{}, err
	}

	// если хотя бы один элемент возможно разместить в остаточном формате
	if remainingLayout.Total() > 0 {
		ri.logger.Debug().MsgFunc(
			func() string {
				return fmt.Sprintf(
					"- placed item %s on remaining out format %s with fict margins: %s, %s, %d * %d = %d",
					item.WithDistance(),
					outRemaining,
					item.WithDistance().OrientationType(),
					remainingPosition,
					remainingLayout.ByWidth,
					remainingLayout.ByHeight,
					remainingLayout.Total(),
				)
			},
		)
	} else {
		ri.logger.Debug().MsgFunc(
			func() string {
				return fmt.Sprintf(
					"- skipped: item format %s not on remaining format %s with fict margins",
					item.WithDistance(),
					outRemaining,
				)
			},
		)
	}

	return remainingLayout, nil
}

func (ri *AlgoRemaining) remainingFormat(item rect.Item, out rect.Format, layout base.Fragment) (format rect.Format, position string) {
	correct := rect.Format{}
	inWithDistance := item.WithDistance()

	if inWithDistance.Width >= inWithDistance.Height {
		position = base.PositionLeft
		format = rect.Format{
			Width:  out.Height,
			Height: out.Width + item.Distance.Width,
		}
		correct.Height = inWithDistance.Width*float64(layout.ByWidth) + item.Distance.Max()
	} else {
		position = base.PositionBottom
		format = rect.Format{
			Width:  out.Height + item.Distance.Height,
			Height: out.Width,
		}
		correct.Width = inWithDistance.Height*float64(layout.ByHeight) + item.Distance.Max()
	}

	// прибавляется фиктивная граница для учёта размещения граничных элементов
	format = format.Sum(item.Distance)

	// вычитается расстояние между вертикально ориентированными
	// элементами и горизонтально ориентированными
	return format.Diff(correct), position
}
