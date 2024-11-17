package adm

import (
	"context"

	"github.com/mondegor/go-storage/mrpostgres/builder"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/adm/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/adm/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/catalog/laminate"
)

func createUnitLaminate(ctx context.Context, opts laminate.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitLaminate(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitLaminate(ctx context.Context, opts laminate.Options) (*httpv1.Laminate, error) {
	entityMeta, err := mrsql.ParseEntity(mrlog.Ctx(ctx), entity.Laminate{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewLaminatePostgres(
		opts.DBConnManager,
		builder.NewSQL(
			builder.WithSQLSetMetaEntity(entityMeta.MetaUpdate()),
			builder.WithSQLOrderByDefaultSort(entityMeta.MetaOrderBy().DefaultSort()),
			builder.WithSQLLimitMaxSize(opts.PageSizeMax),
		),
	)
	useCase := usecase.NewLaminate(storage, opts.MaterialTypeAPI, opts.EventEmitter, opts.UseCaseErrorWrapper)
	controller := httpv1.NewLaminate(
		opts.RequestParsers.ExtendParser,
		opts.ResponseSender,
		useCase,
		entityMeta.MetaOrderBy(),
	)

	return controller, nil
}
