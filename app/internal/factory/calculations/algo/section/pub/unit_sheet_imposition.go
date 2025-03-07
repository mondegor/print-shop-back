package pub

import (
	"context"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/imposition/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/imposition/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/calculations/algo"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/algo/sheet/imposition"
)

func createUnitSheetImposition(ctx context.Context, opts algo.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitSheetImposition(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitSheetImposition(ctx context.Context, opts algo.Options) (*httpv1.SheetImposition, error) { //nolint:unparam
	algoComponent := imposition.New(mrlog.Ctx(ctx))
	useCase := usecase.NewSheetImposition(algoComponent, opts.EventEmitter, opts.UseCaseErrorWrapper)
	controller := httpv1.NewSheetImposition(
		opts.RequestParsers.Validator,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
