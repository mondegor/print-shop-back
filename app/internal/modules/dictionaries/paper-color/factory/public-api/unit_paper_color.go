package factory

import (
	"context"
	http_v1 "print-shop-back/internal/modules/dictionaries/paper-color/controller/http_v1/public-api"
	"print-shop-back/internal/modules/dictionaries/paper-color/factory"
	repository "print-shop-back/internal/modules/dictionaries/paper-color/infrastructure/repository/public-api"
	usecase "print-shop-back/internal/modules/dictionaries/paper-color/usecase/public-api"

	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitPaperColor(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitPaperColor(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitPaperColor(ctx context.Context, opts factory.Options) (*http_v1.PaperColor, error) {
	storage := repository.NewPaperColorPostgres(
		opts.PostgresAdapter,
	)
	useCase := usecase.NewPaperColor(storage, opts.UsecaseHelper)
	controller := http_v1.NewPaperColor(
		opts.RequestParser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
