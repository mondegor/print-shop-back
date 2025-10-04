package pub

import (
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/insideoutside/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/insideoutside/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/calculations/algo"
)

func createUnitSheetInsideOutside(opts algo.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitSheetInsideOutside(opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitSheetInsideOutside(opts algo.Options) (*httpv1.SheetInsideOutside, error) { //nolint:unparam
	useCase := usecase.NewSheetInsideOutside(opts.EventEmitter)
	controller := httpv1.NewSheetInsideOutside(
		opts.RequestParsers.Validator,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
