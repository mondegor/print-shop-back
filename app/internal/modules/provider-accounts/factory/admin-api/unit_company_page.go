package factory

import (
	"context"
	http_v1 "print-shop-back/internal/modules/provider-accounts/controller/http_v1/admin-api"
	entity "print-shop-back/internal/modules/provider-accounts/entity/admin-api"
	"print-shop-back/internal/modules/provider-accounts/factory"
	repository "print-shop-back/internal/modules/provider-accounts/infrastructure/repository/admin-api"
	usecase "print-shop-back/internal/modules/provider-accounts/usecase/admin-api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitCompanyPage(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitCompanyPage(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitCompanyPage(ctx context.Context, opts factory.Options) (*http_v1.CompanyPage, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.CompanyPage{})

	if err != nil {
		return nil, err
	}

	storage := repository.NewCompanyPagePostgres(
		opts.PostgresAdapter,
		mrpostgres.NewSqlBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			mrpostgres.NewSqlBuilderOrderBy(ctx, metaOrderBy.DefaultSort()),
			mrpostgres.NewSqlBuilderPager(opts.PageSizeMax),
		),
	)
	useCase := usecase.NewCompanyPage(storage, opts.UsecaseHelper, opts.UnitCompanyPage.LogoURLBuilder)
	controller := http_v1.NewCompanyPage(
		opts.RequestParsers.Parser,
		opts.ResponseSender,
		useCase,
		metaOrderBy,
	)

	return controller, nil
}
