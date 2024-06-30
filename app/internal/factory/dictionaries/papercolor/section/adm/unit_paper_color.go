package adm

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/adm/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/adm/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/papercolor"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitPaperColor(ctx context.Context, opts papercolor.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitPaperColor(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitPaperColor(ctx context.Context, opts papercolor.Options) (*httpv1.PaperColor, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.PaperColor{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewPaperColorPostgres(
		opts.DBConnManager,
		mrpostgres.NewSQLBuilderSelect(
			mrpostgres.NewSQLBuilderWhere(),
			mrpostgres.NewSQLBuilderOrderBy(ctx, metaOrderBy.DefaultSort()),
			mrpostgres.NewSQLBuilderLimit(opts.PageSizeMax),
		),
	)
	useCase := usecase.NewPaperColor(storage, opts.EventEmitter, opts.UsecaseHelper)
	controller := httpv1.NewPaperColor(
		opts.RequestParsers.ExtendParser,
		opts.ResponseSender,
		useCase,
		metaOrderBy,
	)

	return controller, nil
}
