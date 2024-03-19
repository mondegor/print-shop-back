package factory

import (
	"context"
	http_v1 "print-shop-back/internal/modules/dictionaries/print-format/controller/http_v1/admin-api"
	entity "print-shop-back/internal/modules/dictionaries/print-format/entity/admin-api"
	"print-shop-back/internal/modules/dictionaries/print-format/factory"
	repository "print-shop-back/internal/modules/dictionaries/print-format/infrastructure/repository/admin-api"
	usecase "print-shop-back/internal/modules/dictionaries/print-format/usecase/admin-api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitPrintFormat(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitPrintFormat(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitPrintFormat(ctx context.Context, opts factory.Options) (*http_v1.PrintFormat, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.PrintFormat{})

	if err != nil {
		return nil, err
	}

	storage := repository.NewPrintFormatPostgres(
		opts.PostgresAdapter,
		mrpostgres.NewSqlBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			mrpostgres.NewSqlBuilderOrderBy(ctx, metaOrderBy.DefaultSort()),
			mrpostgres.NewSqlBuilderPager(opts.PageSizeMax),
		),
	)
	useCase := usecase.NewPrintFormat(storage, opts.EventEmitter, opts.UsecaseHelper)
	controller := http_v1.NewPrintFormat(
		opts.RequestParser,
		opts.ResponseSender,
		useCase,
		metaOrderBy,
	)

	return controller, nil
}
