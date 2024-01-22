package factory

import (
	module "print-shop-back/internal/modules/controls"
	http_v1 "print-shop-back/internal/modules/controls/controller/http_v1/admin-api"
	entity "print-shop-back/internal/modules/controls/entity/admin-api"
	"print-shop-back/internal/modules/controls/factory"
	repository "print-shop-back/internal/modules/controls/infrastructure/repository/admin-api"
	usecase "print-shop-back/internal/modules/controls/usecase/admin-api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitElementTemplate(opts *factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitElementTemplate(opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitElementTemplate(opts *factory.Options) (*http_v1.ElementTemplate, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(entity.ElementTemplate{})

	if err != nil {
		return nil, err
	}

	entityMetaUpdate, err := mrsql.NewEntityMetaUpdate(entity.ElementTemplate{})

	if err != nil {
		return nil, err
	}

	storage := repository.NewElementTemplatePostgres(
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
	service := usecase.NewElementTemplate(storage, opts.EventBox, opts.ServiceHelper)
	controller := http_v1.NewElementTemplate(
		opts.RequestParser,
		opts.ResponseSender,
		service,
		metaOrderBy,
	)

	return controller, nil
}
