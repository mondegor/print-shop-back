package pub

import (
	"context"

	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/cutting/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/cutting/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/calculations/algo"
)

func createUnitRectCutting(ctx context.Context, opts algo.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitRectCutting(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitRectCutting(_ context.Context, opts algo.Options) (*httpv1.RectCutting, error) { //nolint:unparam
	useCase := usecase.NewRectCutting(opts.EventEmitter, opts.UsecaseHelper)
	controller := httpv1.NewRectCutting(
		opts.RequestParsers.Validator,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
