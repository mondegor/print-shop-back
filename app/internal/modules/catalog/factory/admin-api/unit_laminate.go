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

func createUnitLaminate(opts *factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitLaminate(opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitLaminate(opts *factory.Options) (*http_v1.Laminate, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(entity.Laminate{})

	if err != nil {
		return nil, err
	}

	entityMetaUpdate, err := mrsql.NewEntityMetaUpdate(entity.Laminate{})

	if err != nil {
		return nil, err
	}

	storage := repository.NewLaminatePostgres(
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
	service := usecase.NewLaminate(storage, opts.LaminateTypeAPI, opts.EventBox, opts.ServiceHelper)
	controller := http_v1.NewLaminate(
		opts.RequestParser,
		opts.ResponseSender,
		service,
		metaOrderBy,
	)

	return controller, nil
}
