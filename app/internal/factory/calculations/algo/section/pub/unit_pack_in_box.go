package pub

import (
	"context"

	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/circulation/packinbox/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/circulation/packinbox/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/calculations/algo"
)

func createUnitCirculationPackInBox(ctx context.Context, opts algo.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitCirculationPackInBox(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitCirculationPackInBox(_ context.Context, opts algo.Options) (*httpv1.PackInBox, error) { //nolint:unparam
	useCase := usecase.NewCirculationPackInBox(opts.EventEmitter, opts.UsecaseHelper)
	controller := httpv1.NewPackInBox(
		opts.RequestParsers.Validator,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
