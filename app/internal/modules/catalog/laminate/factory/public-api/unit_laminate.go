package factory

import (
	"context"
	http_v1 "print-shop-back/internal/modules/catalog/laminate/controller/http_v1/public-api"
	"print-shop-back/internal/modules/catalog/laminate/factory"
	repository "print-shop-back/internal/modules/catalog/laminate/infrastructure/repository/public-api"
	usecase "print-shop-back/internal/modules/catalog/laminate/usecase/public-api"

	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitLaminate(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitLaminate(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitLaminate(ctx context.Context, opts factory.Options) (*http_v1.Laminate, error) {
	storage := repository.NewLaminatePostgres(
		opts.PostgresAdapter,
	)
	useCase := usecase.NewLaminate(storage, opts.UsecaseHelper)
	controller := http_v1.NewLaminate(
		opts.RequestParser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
