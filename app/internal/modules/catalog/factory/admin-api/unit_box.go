package factory

import (
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

func createUnitBox(opts *factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitBox(opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitBox(opts *factory.Options) (*http_v1.Box, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(entity.Box{})

	if err != nil {
		return nil, err
	}

	entityMetaUpdate, err := mrsql.NewEntityMetaUpdate(entity.Box{})

	if err != nil {
		return nil, err
	}

	storage := repository.NewBoxPostgres(
		opts.PostgresAdapter,
		mrsql.NewBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			mrpostgres.NewSqlBuilderOrderByWithDefaultSort(metaOrderBy.DefaultSort()),
			mrpostgres.NewSqlBuilderPager(module.PageSizeMax),
		),
		mrsql.NewBuilderUpdateWithMeta(
			entityMetaUpdate,
			mrpostgres.NewSqlBuilderSet(),
		),
	)
	service := usecase.NewBox(storage, opts.EventBox, opts.ServiceHelper)
	controller := http_v1.NewBox(
		opts.RequestParser,
		opts.ResponseSender,
		service,
		metaOrderBy,
	)

	return controller, nil
}
