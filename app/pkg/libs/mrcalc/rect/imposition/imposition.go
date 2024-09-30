package imposition

import (
	"errors"
	"fmt"

	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/imposition/remaining"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/imposition/total"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/insideoutside"
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

	// AlgoResult - результат вычислений спуска полос.
	AlgoResult struct {
		Layout    rect.Format     // формат листа, куда помещаются все изделия
		Fragments []base.Fragment // раскладка изделий на листе
		Total     uint64          // общее кол-во элементов
		RestArea  float64         // неиспользуемый остаток (m2)
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
func (ri *Algo) Calc(item rect.Item, out rect.Format, opts Options) (AlgoResult, error) {
	if !item.IsValid() {
		return AlgoResult{}, fmt.Errorf("item format is not valid: %s", item.Format)
	}

	if !out.IsValid() {
		return AlgoResult{}, fmt.Errorf("out format is not valid: %s", out)
	}

	ri.logger.Debug().MsgFunc(
		func() string {
			return fmt.Sprintf(
				"Calc(%s) :: Calculate count rect items of %s + %s on out format %s",
				ri.algoType(opts),
				item.Format,
				item.Distance,
				out,
			)
		},
	)

	outWork := out

	if opts.UseMirror {
		// режется лист пополам по вертикали и размещаются на нём элементы
		// таким образом алгоритм добивается четного и зеркального размещения элементов на весь лист
		outWork = rect.Format{
			Width:  (out.Width - item.Distance.Max()) / 2,
			Height: out.Height,
		}
	}

	mainLayout, err := ri.calcMainLayout(item, outWork)
	if err != nil {
		return AlgoResult{}, err
	}

	if mainLayout.Total() == 0 {
		ri.logger.Debug().MsgFunc(
			func() string {
				return fmt.Sprintf(
					"- calculated unsuccessfully: ZERO items of %s on %s",
					item.Format,
					outWork,
				)
			},
		)

		return AlgoResult{}, nil
	}

	const maxFragments = 2
	fragments := make(base.Fragments, 0, maxFragments)
	fragments = append(fragments, mainLayout)

	if opts.AllowRotation {
		remainingLayout, err := ri.remaining.Calc(mainLayout, item, outWork)
		if err != nil {
			return AlgoResult{}, fmt.Errorf("opts.AllowRotation=true: %w", err)
		}

		if remainingLayout.Total() > 0 {
			fragments = append(fragments, remainingLayout)
		}
	}

	if opts.UseMirror {
		fragments[0].ByWidth *= 2

		if len(fragments) > 1 {
			fragments[1].ByHeight *= 2
		}
	}

	totalLayout, restArea := ri.total.Calc(item, out, fragments)

	ri.logger.Debug().MsgFunc(
		func() string {
			return fmt.Sprintf(
				"- calculated successfully: %d items of %s + %s on %s, image: %s, rest %.2fx%.2f ~= %.2f, fragments: %d",
				fragments.Total(),
				item.Format,
				item.Distance,
				out,
				totalLayout,
				item.Format.Width,
				restArea/item.Format.Width,
				restArea,
				len(fragments),
			)
		},
	)

	if restArea < 0 {
		return AlgoResult{}, errors.New("restArea < 0") // TODO: restArea
	}

	return AlgoResult{
		Layout:    totalLayout,
		Fragments: fragments,
		Total:     fragments.Total(),
		RestArea:  restArea,
	}, nil
}

func (ri *Algo) calcMainLayout(item rect.Item, out rect.Format) (base.Fragment, error) {
	// добавляются фиктивные границы к внешнему формату,
	// для того, чтобы поместить граничные элементы, у которых внешние края
	// в реальности короче (т.к. item.Distance выступает в качестве межэлементного расстояния)
	outWithFictBorders := out.Sum(item.Distance)

	layout, err := insideoutside.AlgoQuantity(item.WithDistance(), outWithFictBorders)
	if err != nil {
		return base.Fragment{}, fmt.Errorf("imposition.Algo.calcMainLayout[item=%+v, outWithFictBorders=%s]: %w", item, outWithFictBorders, err)
	}

	ri.logger.Debug().MsgFunc(
		func() string {
			return fmt.Sprintf(
				"- placed item %s on out format %s with fict margins: %s, %s, %d * %d = %d",
				item.WithDistance(),
				out,
				item.WithDistance().OrientationType(),
				base.PositionTop,
				layout.ByWidth,
				layout.ByHeight,
				layout.Total(),
			)
		},
	)

	return layout, nil
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
