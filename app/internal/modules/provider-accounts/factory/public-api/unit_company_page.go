package factory

import (
	http_v1 "print-shop-back/internal/modules/provider-accounts/controller/http_v1/public-api"
	"print-shop-back/internal/modules/provider-accounts/factory"
	repository "print-shop-back/internal/modules/provider-accounts/infrastructure/repository/public-api"
	usecase "print-shop-back/internal/modules/provider-accounts/usecase/public-api"

	"github.com/mondegor/go-webcore/mrcore"
)

func newUnitCompanyPage(
	c *[]mrcore.HttpController,
	opts *factory.Options,
	section mrcore.ClientSection,
) error {
	storage := repository.NewCompanyPagePostgres(opts.PostgresAdapter)
	service := usecase.NewCompanyPage(storage, opts.ServiceHelper)
	*c = append(*c, http_v1.NewCompanyPage(section, service, opts.UnitCompanyPage.LogoURLBuilder))

	return nil
}
