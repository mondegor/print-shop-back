package factory

import (
	"context"
	http_v1 "print-shop-back/internal/modules/catalog/paper/controller/http_v1/admin-api"
	entity "print-shop-back/internal/modules/catalog/paper/entity/admin-api"
	"print-shop-back/internal/modules/catalog/paper/factory"
	repository "print-shop-back/internal/modules/catalog/paper/infrastructure/repository/admin-api"
	usecase "print-shop-back/internal/modules/catalog/paper/usecase/admin-api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitPaper(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitPaper(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitPaper(ctx context.Context, opts factory.Options) (*http_v1.Paper, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.Paper{})

	if err != nil {
		return nil, err
	}

	entityMetaUpdate, err := mrsql.NewEntityMetaUpdate(ctx, entity.Paper{})

	if err != nil {
		return nil, err
	}

	storage := repository.NewPaperPostgres(
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
	useCase := usecase.NewPaper(storage, opts.PaperColorAPI, opts.PaperFactureAPI, opts.EventEmitter, opts.UsecaseHelper)
	controller := http_v1.NewPaper(
		opts.RequestParser,
		opts.ResponseSender,
		useCase,
		metaOrderBy,
	)

	return controller, nil
}
