package pub

import (
	"context"

	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/box/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/calculations/box/section/pub/repository"
	"github.com/mondegor/print-shop-back/internal/calculations/box/section/pub/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/calculations/box"
)

func createUnitCalcResult(ctx context.Context, opts box.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitCalcResult(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitCalcResult(_ context.Context, opts box.Options) (*httpv1.Box, error) { //nolint:unparam
	storage := repository.NewBoxPostgres(
		opts.DBConnManager,
	)
	useCase := usecase.NewBox(storage, opts.EventEmitter, opts.UsecaseHelper)
	controller := httpv1.NewBox(
		opts.RequestParsers.ExtendParser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
