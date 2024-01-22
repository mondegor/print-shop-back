package factory

import (
	module "print-shop-back/internal/modules/dictionaries"
	http_v1 "print-shop-back/internal/modules/dictionaries/controller/http_v1/admin-api"
	entity "print-shop-back/internal/modules/dictionaries/entity/admin-api"
	"print-shop-back/internal/modules/dictionaries/factory"
	repository "print-shop-back/internal/modules/dictionaries/infrastructure/repository/admin-api"
	usecase "print-shop-back/internal/modules/dictionaries/usecase/admin-api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitPrintFormat(opts *factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitPrintFormat(opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitPrintFormat(opts *factory.Options) (*http_v1.PrintFormat, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(entity.PrintFormat{})

	if err != nil {
		return nil, err
	}

	entityMetaUpdate, err := mrsql.NewEntityMetaUpdate(entity.PrintFormat{})

	if err != nil {
		return nil, err
	}

	storage := repository.NewPrintFormatPostgres(
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
	service := usecase.NewPrintFormat(storage, opts.EventBox, opts.ServiceHelper)
	controller := http_v1.NewPrintFormat(
		opts.RequestParser,
		opts.ResponseSender,
		service,
		metaOrderBy,
	)

	return controller, nil
}
