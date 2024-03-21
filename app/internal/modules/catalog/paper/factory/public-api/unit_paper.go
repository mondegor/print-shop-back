package factory

import (
	"context"
	http_v1 "print-shop-back/internal/modules/catalog/paper/controller/http_v1/public-api"
	"print-shop-back/internal/modules/catalog/paper/factory"
	repository "print-shop-back/internal/modules/catalog/paper/infrastructure/repository/public-api"
	usecase "print-shop-back/internal/modules/catalog/paper/usecase/public-api"

	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitPaper(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitPaper(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitPaper(ctx context.Context, opts factory.Options) (*http_v1.Paper, error) {
	storage := repository.NewPaperPostgres(
		opts.PostgresAdapter,
	)
	useCase := usecase.NewPaper(storage, opts.UsecaseHelper)
	controller := http_v1.NewPaper(
		opts.RequestParser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
