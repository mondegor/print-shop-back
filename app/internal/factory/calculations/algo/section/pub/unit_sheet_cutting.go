package pub

import (
	"context"

	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/cutting/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/cutting/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/calculations/algo"
)

func createUnitSheetCutting(ctx context.Context, opts algo.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitSheetCutting(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitSheetCutting(_ context.Context, opts algo.Options) (*httpv1.SheetCutting, error) { //nolint:unparam
	useCase := usecase.NewSheetCutting(opts.EventEmitter, opts.UseCaseErrorWrapper)
	controller := httpv1.NewSheetCutting(
		opts.RequestParsers.Validator,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
