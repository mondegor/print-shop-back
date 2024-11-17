package pub

import (
	"context"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/circulation/packinbox/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/circulation/packinbox/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/calculations/algo"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/packinbox"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/imposition"
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

func newUnitCirculationPackInBox(ctx context.Context, opts algo.Options) (*httpv1.PackInBox, error) { //nolint:unparam
	logger := mrlog.Ctx(ctx)
	impAlgo := imposition.New(logger)
	packInBoxAlgo := packinbox.New(logger, impAlgo)

	useCase := usecase.NewCirculationPackInBox(packInBoxAlgo, opts.EventEmitter, opts.UseCaseErrorWrapper)
	controller := httpv1.NewPackInBox(
		opts.RequestParsers.Validator,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
