package pub

import (
	"context"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/packinstack/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/packinstack/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/calculations/algo"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/algo/sheet/packinstack"
)

func createUnitSheetPackInStack(ctx context.Context, opts algo.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitSheetPackInStack(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitSheetPackInStack(ctx context.Context, opts algo.Options) (*httpv1.PackInStack, error) { //nolint:unparam
	logger := mrlog.Ctx(ctx)
	packInStackAlgoSheet := packinstack.New(logger)

	useCase := usecase.NewSheetPackInStack(packInStackAlgoSheet, opts.EventEmitter, opts.UseCaseErrorWrapper)
	controller := httpv1.NewSheetPackInStack(
		opts.RequestParsers.Validator,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
