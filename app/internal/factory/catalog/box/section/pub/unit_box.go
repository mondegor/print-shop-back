package pub

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/catalog/box/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/catalog/box/section/pub/repository"
	"github.com/mondegor/print-shop-back/internal/catalog/box/section/pub/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/catalog/box"

	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitBox(ctx context.Context, opts box.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitBox(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitBox(_ context.Context, opts box.Options) (*httpv1.Box, error) { //nolint:unparam
	storage := repository.NewBoxPostgres(
		opts.DBConnManager,
	)
	useCase := usecase.NewBox(storage, opts.UsecaseHelper)
	controller := httpv1.NewBox(
		opts.RequestParsers.Parser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
