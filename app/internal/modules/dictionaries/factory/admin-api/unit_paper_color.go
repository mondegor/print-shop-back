package factory

import (
	"context"
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

func createUnitPaperColor(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitPaperColor(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitPaperColor(ctx context.Context, opts factory.Options) (*http_v1.PaperColor, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.PaperColor{})

	if err != nil {
		return nil, err
	}

	storage := repository.NewPaperColorPostgres(
		opts.PostgresAdapter,
		mrsql.NewBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			mrpostgres.NewSqlBuilderOrderByWithDefaultSort(ctx, metaOrderBy.DefaultSort()),
			mrpostgres.NewSqlBuilderPager(module.PageSizeMax),
		),
	)
	service := usecase.NewPaperColor(storage, opts.EventEmitter, opts.UsecaseHelper)
	controller := http_v1.NewPaperColor(
		opts.RequestParser,
		opts.ResponseSender,
		service,
		metaOrderBy,
	)

	return controller, nil
}
