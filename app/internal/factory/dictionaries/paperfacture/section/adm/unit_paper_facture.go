package adm

import (
	"context"

	"github.com/mondegor/go-storage/mrpostgres/builder"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/adm/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/adm/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/paperfacture"
)

func createUnitPaperFacture(ctx context.Context, opts paperfacture.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitPaperFacture(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitPaperFacture(ctx context.Context, opts paperfacture.Options) (*httpv1.PaperFacture, error) {
	entityMeta, err := mrsql.ParseEntity(mrlog.Ctx(ctx), entity.PaperFacture{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewPaperFacturePostgres(
		opts.DBConnManager,
		builder.NewSQL(
			builder.WithSQLOrderByDefaultSort(entityMeta.MetaOrderBy().DefaultSort()),
			builder.WithSQLLimitMaxSize(opts.PageSizeMax),
		),
	)
	useCase := usecase.NewPaperFacture(storage, opts.EventEmitter, opts.UseCaseErrorWrapper)
	controller := httpv1.NewPaperFacture(
		opts.RequestParsers.ExtendParser,
		opts.ResponseSender,
		useCase,
		entityMeta.MetaOrderBy(),
	)

	return controller, nil
}
