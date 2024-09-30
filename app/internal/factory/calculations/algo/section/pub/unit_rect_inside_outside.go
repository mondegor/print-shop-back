package pub

import (
	"context"

	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/insideoutside/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/insideoutside/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/calculations/algo"
)

func createUnitRectInsideOutside(ctx context.Context, opts algo.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitRectInsideOutside(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitRectInsideOutside(_ context.Context, opts algo.Options) (*httpv1.RectInsideOutside, error) { //nolint:unparam
	useCase := usecase.NewRectInsideOutside(opts.EventEmitter, opts.UseCaseHelper)
	controller := httpv1.NewRectInsideOutside(
		opts.RequestParsers.Validator,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
