package factory

import (
	"context"
	module "print-shop-back/internal/modules/dictionaries/laminate-type"
	http_v1 "print-shop-back/internal/modules/dictionaries/laminate-type/controller/http_v1/admin-api"
	entity "print-shop-back/internal/modules/dictionaries/laminate-type/entity/admin-api"
	"print-shop-back/internal/modules/dictionaries/laminate-type/factory"
	repository "print-shop-back/internal/modules/dictionaries/laminate-type/infrastructure/repository/admin-api"
	usecase "print-shop-back/internal/modules/dictionaries/laminate-type/usecase/admin-api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitLaminateType(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitLaminateType(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitLaminateType(ctx context.Context, opts factory.Options) (*http_v1.LaminateType, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.LaminateType{})

	if err != nil {
		return nil, err
	}

	storage := repository.NewLaminateTypePostgres(
		opts.PostgresAdapter,
		mrsql.NewBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			mrpostgres.NewSqlBuilderOrderByWithDefaultSort(ctx, metaOrderBy.DefaultSort()),
			mrpostgres.NewSqlBuilderPager(module.PageSizeMax),
		),
	)
	service := usecase.NewLaminateType(storage, opts.EventEmitter, opts.UsecaseHelper)
	controller := http_v1.NewLaminateType(
		opts.RequestParser,
		opts.ResponseSender,
		service,
		metaOrderBy,
	)

	return controller, nil
}
