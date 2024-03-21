package factory

import (
	"context"
	http_v1 "print-shop-back/internal/modules/dictionaries/paper-facture/controller/http_v1/public-api"
	"print-shop-back/internal/modules/dictionaries/paper-facture/factory"
	repository "print-shop-back/internal/modules/dictionaries/paper-facture/infrastructure/repository/public-api"
	usecase "print-shop-back/internal/modules/dictionaries/paper-facture/usecase/public-api"

	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitPaperFacture(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitPaperFacture(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitPaperFacture(ctx context.Context, opts factory.Options) (*http_v1.PaperFacture, error) {
	storage := repository.NewPaperFacturePostgres(
		opts.PostgresAdapter,
	)
	useCase := usecase.NewPaperFacture(storage, opts.UsecaseHelper)
	controller := http_v1.NewPaperFacture(
		opts.RequestParser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
