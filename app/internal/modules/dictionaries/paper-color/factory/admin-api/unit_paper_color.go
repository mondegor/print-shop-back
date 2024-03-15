package factory

import (
	"context"
	http_v1 "print-shop-back/internal/modules/dictionaries/paper-color/controller/http_v1/admin-api"
	entity "print-shop-back/internal/modules/dictionaries/paper-color/entity/admin-api"
	"print-shop-back/internal/modules/dictionaries/paper-color/factory"
	repository "print-shop-back/internal/modules/dictionaries/paper-color/infrastructure/repository/admin-api"
	usecase "print-shop-back/internal/modules/dictionaries/paper-color/usecase/admin-api"

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
		mrpostgres.NewSqlBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			mrpostgres.NewSqlBuilderOrderBy(ctx, metaOrderBy.DefaultSort()),
			mrpostgres.NewSqlBuilderPager(opts.PageSizeMax),
		),
	)
	useCase := usecase.NewPaperColor(storage, opts.EventEmitter, opts.UsecaseHelper)
	controller := http_v1.NewPaperColor(
		opts.RequestParser,
		opts.ResponseSender,
		useCase,
		metaOrderBy,
	)

	return controller, nil
}
