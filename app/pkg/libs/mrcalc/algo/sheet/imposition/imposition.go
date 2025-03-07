package imposition

import (
	"errors"
	"fmt"

	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/algo/sheet/imposition/remaining"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/algo/sheet/imposition/total"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/algo/sheet/insideoutside"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/enum"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2/rect2d"
)

type (
	// Algo - спуск полос рассчитывает раскладку максимального количество элементов
	// указанного формата, которое можно разместить на листе указанного формата.
	// Также учитывается расстояние между элементами по горизонтали и вертикали.
	Algo struct {
		logger    mrlog.Logger
		remaining *remaining.AlgoRemaining
		total     *total.AlgoTotal
	}

	// Options - опции алгоритма Algo.
	Options struct {
		AllowRotation bool // при true разрешается располагать элементы повёрнутые на 90 градусов друг к другу
		UseMirror     bool // размещение на листе элементов зеркально друг к другу
	}

	// Output - результат вычислений спуска полос.
	Output struct {
		ContainerFormat rect2d.Format    // минимальный достаточный формат листа, куда помещаются все изделия
		Fragments       rect2d.Fragments // раскладка изделий на листе
		RestArea        float64          // неиспользуемый остаток (m2)
		AllowRotation   bool
		UseMirror       bool
	}
)

// New - создаёт объект Algo.
func New(logger mrlog.Logger) *Algo {
	return &Algo{
		logger:    logger,
		remaining: remaining.New(logger),
		total:     total.New(),
	}
}

// Calc - расчёт алгоритма.
func (ri *Algo) Calc(element, distance, out rect2d.Format, opts Options) (Output, error) {
	if !element.IsValid() {
		return Output{}, fmt.Errorf("element format is not valid: %s", element)
	}

	if !distance.IsZero() && !distance.IsValid() {
		return Output{}, fmt.Errorf("distance format is not valid: %s", distance)
	}

	if !out.IsValid() {
		return Output{}, fmt.Errorf("out format is not valid: %s", out)
	}

	ri.logger.Debug().MsgFunc(
		func() string {
			return fmt.Sprintf(
				"Calc(%s) :: Calculate count rect elements of %s + %s on out format %s",
				ri.algoType(opts),
				element,
				distance,
				out,
			)
		},
	)

	outWork := out

	if opts.UseMirror {
		// режется лист пополам по вертикали и размещаются на нём элементы
		// таким образом алгоритм добивается четного и зеркального размещения элементов на весь лист
		outWork = rect2d.Format{
			Width:  (out.Width - distance.Max()) / 2,
			Height: out.Height,
		}
	}

	topFragment, err := ri.calcTopFragment(element, distance, outWork)
	if err != nil {
		return Output{}, err
	}

	if topFragment.Layout.Quantity() == 0 {
		ri.logger.Debug().MsgFunc(
			func() string {
				return fmt.Sprintf(
					"- calculated unsuccessfully: ZERO elements of %s on %s",
					topFragment.Element,
					outWork,
				)
			},
		)

		return Output{}, nil
	}

	const maxFragments = 2
	fragments := make(rect2d.Fragments, 1, maxFragments)
	fragments[0] = topFragment

	if opts.AllowRotation {
		remainingFragment, err := ri.remaining.Calc(topFragment, outWork)
		if err != nil {
			return Output{}, fmt.Errorf("opts.AllowRotation=true: %w", err)
		}

		if remainingFragment.Layout.Quantity() > 0 {
			fragments = append(fragments, remainingFragment)
		}
	}

	if opts.UseMirror {
		fragments[0].Layout.ByWidth *= 2

		if len(fragments) > 1 {
			fragments[1].Layout.ByHeight *= 2
		}
	}

	containerFormat, restArea := ri.total.Calc(fragments, out)

	ri.logger.Debug().MsgFunc(
		func() string {
			return fmt.Sprintf(
				"- calculated successfully: %d elements of %s + %s on %s, image: %s, rest %.2fx%.2f ~= %.2f, fragments: %d",
				fragments.TotalQuantity(),
				element,
				distance,
				out,
				containerFormat,
				element.Width,
				restArea/element.Width,
				restArea,
				len(fragments),
			)
		},
	)

	if restArea < 0 {
		return Output{}, errors.New("restArea < 0") // TODO: restArea
	}

	return Output{
		ContainerFormat: containerFormat,
		Fragments:       fragments,
		RestArea:        restArea,
		AllowRotation:   opts.AllowRotation,
		UseMirror:       opts.UseMirror,
	}, nil
}

func (ri *Algo) calcTopFragment(element, distance, out rect2d.Format) (rect2d.Fragment, error) {
	// добавляются фиктивные границы к внешнему формату,
	// для того, чтобы поместить граничные элементы, у которых внешние края
	// в реальности короче (т.к. distance выступает в качестве межэлементного расстояния)
	outWithFictBorders := out.Add(distance)

	layout, err := insideoutside.AlgoQuantity(element.Add(distance), outWithFictBorders)
	if err != nil {
		return rect2d.Fragment{}, fmt.Errorf(
			"imposition.Algo.calcTopFragment[elementWithDistance=%+v, outWithFictBorders=%s]: %w",
			element.Add(distance),
			outWithFictBorders,
			err,
		)
	}

	topFragment := rect2d.Fragment{
		Element:  element,
		Distance: distance,
		Position: enum.PositionTop,
		Layout:   layout,
	}

	ri.logger.Debug().MsgFunc(
		func() string {
			return fmt.Sprintf(
				"- placed element %s on out format %s with fict margins: %s, %s, %d * %d = %d",
				topFragment.ElementWithDistance(),
				out,
				topFragment.ElementWithDistance().OrientationType(),
				topFragment.Position,
				topFragment.Layout.ByWidth,
				topFragment.Layout.ByHeight,
				topFragment.Layout.Quantity(),
			)
		},
	)

	return topFragment, nil
}

func (ri *Algo) algoType(opts Options) string {
	t := "WITH ROTATION"

	if !opts.AllowRotation {
		t = "ROTATION OFF"
	}

	if opts.UseMirror {
		t += ", MIRROR"
	}

	return t
}
