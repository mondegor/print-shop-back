package factory

import (
	"context"
	http_v1 "print-shop-back/internal/modules/provider-accounts/controller/http_v1/provider-account-api"
	"print-shop-back/internal/modules/provider-accounts/factory"
	repository "print-shop-back/internal/modules/provider-accounts/infrastructure/repository/provider-account-api"
	usecase "print-shop-back/internal/modules/provider-accounts/usecase/provider-account-api"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresponse"
)

func createUnitCompanyPage(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitCompanyPage(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	if c, err := newUnitCompanyPageLogo(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitCompanyPage(ctx context.Context, opts factory.Options) (*http_v1.CompanyPage, error) {
	storage := repository.NewCompanyPagePostgres(opts.PostgresAdapter)
	useCase := usecase.NewCompanyPage(storage, opts.EventEmitter, opts.UsecaseHelper, opts.UnitCompanyPage.LogoURLBuilder)
	controller := http_v1.NewCompanyPage(
		opts.RequestParsers.Parser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}

func newUnitCompanyPageLogo(ctx context.Context, opts factory.Options) (*http_v1.CompanyPageLogo, error) {
	storage := repository.NewCompanyPageLogoPostgres(opts.PostgresAdapter)
	useCase := usecase.NewCompanyPageLogo(
		storage,
		opts.UnitCompanyPage.LogoFileAPI,
		opts.Locker,
		opts.EventEmitter,
		opts.UsecaseHelper,
	)
	controller := http_v1.NewCompanyPageLogo(
		opts.RequestParsers.Image,
		mrresponse.NewFileSender(opts.ResponseSender),
		useCase,
	)

	return controller, nil
}
