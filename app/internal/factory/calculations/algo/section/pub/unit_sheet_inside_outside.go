package pub

import (
	"context"

	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/insideoutside/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/insideoutside/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/calculations/algo"
)

func createUnitSheetInsideOutside(ctx context.Context, opts algo.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitSheetInsideOutside(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitSheetInsideOutside(_ context.Context, opts algo.Options) (*httpv1.SheetInsideOutside, error) { //nolint:unparam
	useCase := usecase.NewSheetInsideOutside(opts.EventEmitter, opts.UseCaseErrorWrapper)
	controller := httpv1.NewSheetInsideOutside(
		opts.RequestParsers.Validator,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
