package factory

import (
	"context"
	http_v1 "print-shop-back/internal/modules/catalog/box/controller/http_v1/admin-api"
	entity "print-shop-back/internal/modules/catalog/box/entity/admin-api"
	"print-shop-back/internal/modules/catalog/box/factory"
	repository "print-shop-back/internal/modules/catalog/box/infrastructure/repository/admin-api"
	usecase "print-shop-back/internal/modules/catalog/box/usecase/admin-api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitBox(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitBox(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitBox(ctx context.Context, opts factory.Options) (*http_v1.Box, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.Box{})

	if err != nil {
		return nil, err
	}

	entityMetaUpdate, err := mrsql.NewEntityMetaUpdate(ctx, entity.Box{})

	if err != nil {
		return nil, err
	}

	storage := repository.NewBoxPostgres(
		opts.PostgresAdapter,
		mrpostgres.NewSqlBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			mrpostgres.NewSqlBuilderOrderBy(ctx, metaOrderBy.DefaultSort()),
			mrpostgres.NewSqlBuilderPager(opts.PageSizeMax),
		),
		mrpostgres.NewSqlBuilderUpdateWithMeta(
			entityMetaUpdate,
			mrpostgres.NewSqlBuilderSet(),
			nil,
		),
	)
	useCase := usecase.NewBox(storage, opts.EventEmitter, opts.UsecaseHelper)
	controller := http_v1.NewBox(
		opts.RequestParser,
		opts.ResponseSender,
		useCase,
		metaOrderBy,
	)

	return controller, nil
}
