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
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(entity.CompanyPage{})

	if err != nil {
		return nil, err
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
	controller := http_v1.NewCompanyPage(
		opts.RequestParsers.Parser,
		opts.ResponseSender,
		service,
		metaOrderBy,
	)

	return controller, nil
}
