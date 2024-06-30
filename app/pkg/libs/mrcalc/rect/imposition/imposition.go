package imposition

import (
	"errors"
	"fmt"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/imposition/remaining"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/imposition/total"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/insideoutside"

	"github.com/mondegor/go-webcore/mrlog"
)

type (
	// Algo - спуск полос.
	// * Рассчитывается максимальное количество элементов указанного формата,
	// * которое можно разместить на листе указанного формата.
	// * Также учитывается расстояние между элементами по горизонтали и вертикали.
	Algo struct {
		logger    mrlog.Logger
		remaining *remaining.AlgoRemaining
		total     *total.AlgoTotal
	}

	// Options - comment struct.
	Options struct {
		AllowRotation bool
		UseMirror     bool
	}

	// AlgoResult - результат вычислений спуска полос.
	AlgoResult struct {
		Layout    rect.Format     `json:"layout"`    // формат листа, куда помещаются все изделия
		Fragments []base.Fragment `json:"fragments"` // раскладка изделий на листе
		Total     uint64          `json:"total"`     // общее кол-во элементов
		RestArea  uint64          `json:"restArea"`  // неиспользуемый остаток (mm2)
	}
)

// New - создаёт объект Algo.
// Поддерживается параметр allowRotation при true разрешается
// располагать элементы повёрнутые на 90 градусов к друг другу.
func New(logger mrlog.Logger) *Algo {
	return &Algo{
		logger:    logger,
		remaining: remaining.New(logger),
		total:     total.New(),
	}
}

// Calc - comment method.
func (ri *Algo) Calc(item rect.Item, out rect.Format, opts Options) (AlgoResult, error) {
	if !item.IsValid() {
		return AlgoResult{}, fmt.Errorf("item format is incorrect: %s", item.Format.String())
	}

	if !out.IsValid() {
		return AlgoResult{}, fmt.Errorf("out format is incorrect: %s", out.String())
	}

	ri.logger.Debug().MsgFunc(
		func() string {
			return fmt.Sprintf(
				"Calc(%s) :: Calculate count rect items of %s + %s on out format %s",
				ri.algoType(opts),
				item.Format,
				item.Border,
				out,
			)
		},
	)

	outWork := out

	if opts.UseMirror {
		// режется лист пополам по вертикали и размещаются на нём элементы
		// таким образом алгоритм добивается четного и зеркального размещения элементов на весь лист
		outWork = rect.Format{
			Width:  (out.Width - item.Border.Max()) / 2,
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
			return AlgoResult{}, err
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
				"- calculated successfully: %d items of %s + %s on %s, image: %s, rest %dx%.2f ~= %d, fragments: %d",
				fragments.Total(),
				item.Format,
				item.Border,
				out,
				totalLayout,
				item.Format.Width,
				float64(restArea/item.Format.Width),
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
		RestArea:  uint64(restArea),
	}, nil
}

// Calc - comment method.
func (ri *Algo) calcMainLayout(item rect.Item, out rect.Format) (base.Fragment, error) {
	// добавляются фиктивные границы к внешнему формату,
	// для того, чтобы поместить граничные элементы, у которых внешние края
	// в реальности короче (т.к. item.Border выступает в качестве межэлементного расстояния)
	outWithBorder := out.Sum(item.Border)

	layout, err := insideoutside.AlgoQuantity(item.WithBorder(), outWithBorder)
	if err != nil {
		return base.Fragment{}, err
	}

	ri.logger.Debug().MsgFunc(
		func() string {
			inWithBorder := item.WithBorder()

			return fmt.Sprintf(
				"- placed item %s on out format %s with fict borders: %s, %s, %d * %d = %d",
				inWithBorder,
				out,
				inWithBorder.OrientationType(),
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
