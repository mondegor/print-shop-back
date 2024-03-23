package factory

import (
	"context"
	http_v1 "print-shop-back/internal/modules/controls/submit-form/controller/http_v1/public-api"
	"print-shop-back/internal/modules/controls/submit-form/factory"
	repository "print-shop-back/internal/modules/controls/submit-form/infrastructure/repository/public-api"
	usecase "print-shop-back/internal/modules/controls/submit-form/usecase/public-api"

	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitSubmitForm(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitSubmitForm(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitSubmitForm(ctx context.Context, opts factory.Options) (*http_v1.SubmitForm, error) {
	storage := repository.NewSubmitFormPostgres(
		opts.PostgresAdapter,
	)
	useCase := usecase.NewSubmitForm(storage, opts.UsecaseHelper)
	controller := http_v1.NewSubmitForm(
		opts.RequestParser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
