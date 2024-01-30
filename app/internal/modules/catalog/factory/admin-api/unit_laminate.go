package factory

import (
	"context"
	module "print-shop-back/internal/modules/catalog"
	http_v1 "print-shop-back/internal/modules/catalog/controller/http_v1/admin-api"
	entity "print-shop-back/internal/modules/catalog/entity/admin-api"
	"print-shop-back/internal/modules/catalog/factory"
	repository "print-shop-back/internal/modules/catalog/infrastructure/repository/admin-api"
	usecase "print-shop-back/internal/modules/catalog/usecase/admin-api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitLaminate(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitLaminate(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitLaminate(ctx context.Context, opts factory.Options) (*http_v1.Laminate, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.Laminate{})

	if err != nil {
		return nil, err
	}

	entityMetaUpdate, err := mrsql.NewEntityMetaUpdate(ctx, entity.Laminate{})

	if err != nil {
		return nil, err
	}

	storage := repository.NewLaminatePostgres(
		opts.PostgresAdapter,
		mrsql.NewBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			mrpostgres.NewSqlBuilderOrderByWithDefaultSort(ctx, metaOrderBy.DefaultSort()),
			mrpostgres.NewSqlBuilderPager(module.PageSizeMax),
		),
		mrsql.NewBuilderUpdateWithMeta(
			entityMetaUpdate,
			mrpostgres.NewSqlBuilderSet(),
		),
	)
	service := usecase.NewLaminate(storage, opts.LaminateTypeAPI, opts.EventEmitter, opts.UsecaseHelper)
	controller := http_v1.NewLaminate(
		opts.RequestParser,
		opts.ResponseSender,
		service,
		metaOrderBy,
	)

	return controller, nil
}
