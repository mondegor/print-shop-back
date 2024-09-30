package pub

import (
	"context"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/imposition/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/imposition/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/calculations/algo"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/imposition"
)

func createUnitRectImposition(ctx context.Context, opts algo.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitRectImposition(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitRectImposition(ctx context.Context, opts algo.Options) (*httpv1.RectImposition, error) { //nolint:unparam
	algoComponent := imposition.New(mrlog.Ctx(ctx))
	useCase := usecase.NewRectImposition(algoComponent, opts.EventEmitter, opts.UseCaseHelper)
	controller := httpv1.NewRectImposition(
		opts.RequestParsers.Validator,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
