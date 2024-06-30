package adm

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/adm/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/adm/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/printformat"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitPrintFormat(ctx context.Context, opts printformat.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitPrintFormat(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitPrintFormat(ctx context.Context, opts printformat.Options) (*httpv1.PrintFormat, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.PrintFormat{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewPrintFormatPostgres(
		opts.DBConnManager,
		mrpostgres.NewSQLBuilderSelect(
			mrpostgres.NewSQLBuilderWhere(),
			mrpostgres.NewSQLBuilderOrderBy(ctx, metaOrderBy.DefaultSort()),
			mrpostgres.NewSQLBuilderLimit(opts.PageSizeMax),
		),
	)
	useCase := usecase.NewPrintFormat(storage, opts.EventEmitter, opts.UsecaseHelper)
	controller := httpv1.NewPrintFormat(
		opts.RequestParsers.ExtendParser,
		opts.ResponseSender,
		useCase,
		metaOrderBy,
	)

	return controller, nil
}
