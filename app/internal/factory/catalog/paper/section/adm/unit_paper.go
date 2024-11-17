package adm

import (
	"context"

	"github.com/mondegor/go-storage/mrpostgres/builder"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/adm/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/adm/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/catalog/paper"
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
	entityMeta, err := mrsql.ParseEntity(mrlog.Ctx(ctx), entity.Paper{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewPaperPostgres(
		opts.DBConnManager,
		builder.NewSQL(
			builder.WithSQLSetMetaEntity(entityMeta.MetaUpdate()),
			builder.WithSQLOrderByDefaultSort(entityMeta.MetaOrderBy().DefaultSort()),
			builder.WithSQLLimitMaxSize(opts.PageSizeMax),
		),
	)
	useCase := usecase.NewPaper(
		storage,
		opts.MaterialTypeAPI,
		opts.PaperColorAPI,
		opts.PaperFactureAPI,
		opts.EventEmitter,
		opts.UseCaseErrorWrapper,
	)
	controller := httpv1.NewPaper(
		opts.RequestParsers.ExtendParser,
		opts.ResponseSender,
		useCase,
		entityMeta.MetaOrderBy(),
	)

	return controller, nil
}
