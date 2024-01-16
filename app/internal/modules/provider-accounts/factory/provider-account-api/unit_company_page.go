package factory

import (
	http_v1 "print-shop-back/internal/modules/provider-accounts/controller/http_v1/provider-account-api"
	"print-shop-back/internal/modules/provider-accounts/factory"
	repository "print-shop-back/internal/modules/provider-accounts/infrastructure/repository/provider-account-api"
	usecase "print-shop-back/internal/modules/provider-accounts/usecase/provider-account-api"

	"github.com/mondegor/go-webcore/mrcore"
)

func newUnitCompanyPage(
	c *[]mrcore.HttpController,
	opts *factory.Options,
	section mrcore.ClientSection,
) error {
	if err := newUnitCompanyPageLogo(c, opts, section); err != nil {
		return err
	}

	return newUnitCompanyPageMain(c, opts, section)
}

func newUnitCompanyPageMain(
	c *[]mrcore.HttpController,
	opts *factory.Options,
	section mrcore.ClientSection,
) error {
	storage := repository.NewCompanyPagePostgres(opts.PostgresAdapter)
	service := usecase.NewCompanyPage(storage, opts.EventBox, opts.ServiceHelper, opts.UnitCompanyPage.LogoURLBuilder)
	*c = append(*c, http_v1.NewCompanyPage(section, service))

	return nil
}

func newUnitCompanyPageLogo(
	c *[]mrcore.HttpController,
	opts *factory.Options,
	section mrcore.ClientSection,
) error {
	storage := repository.NewCompanyPageLogoPostgres(opts.PostgresAdapter)
	service := usecase.NewCompanyPageLogo(
		storage,
		opts.UnitCompanyPage.LogoFileAPI,
		opts.Locker,
		opts.EventBox,
		opts.ServiceHelper,
	)
	*c = append(*c, http_v1.NewCompanyPageLogo(section, service))

	return nil
}
