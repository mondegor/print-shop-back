package adm

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/factory/provideraccounts"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/adm/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/adm/usecase"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitCompanyPage(ctx context.Context, opts provideraccounts.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitCompanyPage(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitCompanyPage(ctx context.Context, opts provideraccounts.Options) (*httpv1.CompanyPage, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.CompanyPage{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewCompanyPagePostgres(
		opts.DBConnManager,
		mrpostgres.NewSQLBuilderSelect(
			mrpostgres.NewSQLBuilderWhere(),
			mrpostgres.NewSQLBuilderOrderBy(ctx, metaOrderBy.DefaultSort()),
			mrpostgres.NewSQLBuilderLimit(opts.PageSizeMax),
		),
	)
	useCase := usecase.NewCompanyPage(storage, opts.UsecaseHelper, opts.UnitCompanyPage.LogoURLBuilder)
	controller := httpv1.NewCompanyPage(
		opts.RequestParsers.ModuleParser,
		opts.ResponseSender,
		useCase,
		metaOrderBy,
	)

	return controller, nil
}
