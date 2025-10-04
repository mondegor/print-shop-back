package pub

import (
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/packinstack/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/packinstack/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/calculations/algo"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/algo/sheet/packinstack"
)

func createUnitSheetPackInStack(opts algo.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitSheetPackInStack(opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitSheetPackInStack(opts algo.Options) (*httpv1.PackInStack, error) { //nolint:unparam
	packInStackAlgoSheet := packinstack.New()

	useCase := usecase.NewSheetPackInStack(packInStackAlgoSheet, opts.EventEmitter)
	controller := httpv1.NewSheetPackInStack(
		opts.RequestParsers.Validator,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
