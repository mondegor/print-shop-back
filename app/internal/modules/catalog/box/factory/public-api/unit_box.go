package factory

import (
	"context"
	http_v1 "print-shop-back/internal/modules/catalog/box/controller/http_v1/public-api"
	"print-shop-back/internal/modules/catalog/box/factory"
	repository "print-shop-back/internal/modules/catalog/box/infrastructure/repository/public-api"
	usecase "print-shop-back/internal/modules/catalog/box/usecase/public-api"

	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitBox(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitBox(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitBox(ctx context.Context, opts factory.Options) (*http_v1.Box, error) {
	storage := repository.NewBoxPostgres(
		opts.PostgresAdapter,
	)
	useCase := usecase.NewBox(storage, opts.UsecaseHelper)
	controller := http_v1.NewBox(
		opts.RequestParser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
