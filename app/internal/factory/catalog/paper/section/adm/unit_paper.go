package adm

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/adm/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/adm/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/catalog/paper"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitPaper(ctx context.Context, opts paper.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitPaper(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitPaper(ctx context.Context, opts paper.Options) (*httpv1.Paper, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.Paper{})
	if err != nil {
		return nil, err
	}

	entityMetaUpdate, err := mrsql.NewEntityMetaUpdate(ctx, entity.Paper{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewPaperPostgres(
		opts.DBConnManager,
		mrpostgres.NewSQLBuilderSelect(
			mrpostgres.NewSQLBuilderWhere(),
			mrpostgres.NewSQLBuilderOrderBy(ctx, metaOrderBy.DefaultSort()),
			mrpostgres.NewSQLBuilderLimit(opts.PageSizeMax),
		),
		mrpostgres.NewSQLBuilderUpdateWithMeta(
			entityMetaUpdate,
			mrpostgres.NewSQLBuilderSet(),
			nil,
		),
	)
	useCase := usecase.NewPaper(
		storage,
		opts.MaterialTypeAPI,
		opts.PaperColorAPI,
		opts.PaperFactureAPI,
		opts.EventEmitter,
		opts.UsecaseHelper,
	)
	controller := httpv1.NewPaper(
		opts.RequestParsers.ExtendParser,
		opts.ResponseSender,
		useCase,
		metaOrderBy,
	)

	return controller, nil
}
