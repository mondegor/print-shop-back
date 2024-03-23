package factory

import (
	"context"
	"print-shop-back/internal/modules/controls/submit-form/factory"
	repository "print-shop-back/internal/modules/controls/submit-form/infrastructure/repository/admin-api"
)

func initUnitSubmitFormVersionEnvironment(ctx context.Context, opts factory.Options) (formVersionOptions, error) {
	storage := repository.NewFormVersionPostgres(
		opts.PostgresAdapter,
	)

	return formVersionOptions{
		storage: storage,
	}, nil
}
