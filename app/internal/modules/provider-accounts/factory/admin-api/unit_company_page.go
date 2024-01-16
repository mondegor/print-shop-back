package factory

import (
	module "print-shop-back/internal/modules/provider-accounts"
	http_v1 "print-shop-back/internal/modules/provider-accounts/controller/http_v1/admin-api"
	entity "print-shop-back/internal/modules/provider-accounts/entity/admin-api"
	"print-shop-back/internal/modules/provider-accounts/factory"
	repository "print-shop-back/internal/modules/provider-accounts/infrastructure/repository/admin-api"
	usecase "print-shop-back/internal/modules/provider-accounts/usecase/admin-api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrcore"
)

func newUnitCompanyPage(
	c *[]mrcore.HttpController,
	opts *factory.Options,
	section mrcore.ClientSection,
) error {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(entity.CompanyPage{})

	if err != nil {
		return err
	}

	storage := repository.NewCompanyPagePostgres(
		opts.PostgresAdapter,
		mrsql.NewBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			mrpostgres.NewSqlBuilderOrderByWithDefaultSort(metaOrderBy.DefaultSort()),
			mrpostgres.NewSqlBuilderPager(module.PageSizeMax),
		),
	)
	service := usecase.NewCompanyPage(storage, opts.ServiceHelper, opts.UnitCompanyPage.LogoURLBuilder)
	*c = append(*c, http_v1.NewCompanyPage(section, service, metaOrderBy))

	return nil
}
