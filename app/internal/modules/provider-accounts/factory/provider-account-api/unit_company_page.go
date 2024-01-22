package factory

import (
	http_v1 "print-shop-back/internal/modules/provider-accounts/controller/http_v1/provider-account-api"
	"print-shop-back/internal/modules/provider-accounts/factory"
	repository "print-shop-back/internal/modules/provider-accounts/infrastructure/repository/provider-account-api"
	usecase "print-shop-back/internal/modules/provider-accounts/usecase/provider-account-api"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresponse"
)

func createUnitCompanyPage(opts *factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitCompanyPage(opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	if c, err := newUnitCompanyPageLogo(opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitCompanyPage(opts *factory.Options) (*http_v1.CompanyPage, error) {
	storage := repository.NewCompanyPagePostgres(opts.PostgresAdapter)
	service := usecase.NewCompanyPage(storage, opts.EventBox, opts.ServiceHelper, opts.UnitCompanyPage.LogoURLBuilder)
	controller := http_v1.NewCompanyPage(
		opts.RequestParsers.Parser,
		opts.ResponseSender,
		service,
	)

	return controller, nil
}

func newUnitCompanyPageLogo(opts *factory.Options) (*http_v1.CompanyPageLogo, error) {
	storage := repository.NewCompanyPageLogoPostgres(opts.PostgresAdapter)
	service := usecase.NewCompanyPageLogo(
		storage,
		opts.UnitCompanyPage.LogoFileAPI,
		opts.Locker,
		opts.EventBox,
		opts.ServiceHelper,
	)
	controller := http_v1.NewCompanyPageLogo(
		mrresponse.NewFileSender(opts.ResponseSender),
		service,
	)

	return controller, nil
}
