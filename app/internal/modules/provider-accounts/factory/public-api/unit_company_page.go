package factory

import (
	http_v1 "print-shop-back/internal/modules/provider-accounts/controller/http_v1/public-api"
	"print-shop-back/internal/modules/provider-accounts/factory"
	repository "print-shop-back/internal/modules/provider-accounts/infrastructure/repository/public-api"
	usecase "print-shop-back/internal/modules/provider-accounts/usecase/public-api"

	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitCompanyPage(opts *factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitCompanyPage(opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitCompanyPage(opts *factory.Options) (*http_v1.CompanyPage, error) {
	storage := repository.NewCompanyPagePostgres(opts.PostgresAdapter)
	service := usecase.NewCompanyPage(storage, opts.ServiceHelper)
	controller := http_v1.NewCompanyPage(
		opts.RequestParsers.String,
		opts.ResponseSender,
		service,
		opts.UnitCompanyPage.LogoURLBuilder,
	)

	return controller, nil
}
