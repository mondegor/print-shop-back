package adm

import (
	"context"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/catalog/box/section/adm/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/catalog/box/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/catalog/box/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/catalog/box/section/adm/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/catalog/box"
)

func createUnitBox(ctx context.Context, opts box.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitBox(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitBox(ctx context.Context, opts box.Options) (*httpv1.Box, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.Box{})
	if err != nil {
		return nil, err
	}

	entityMetaUpdate, err := mrsql.NewEntityMetaUpdate(ctx, entity.Box{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewBoxPostgres(
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
	useCase := usecase.NewBox(storage, opts.EventEmitter, opts.UseCaseHelper)
	controller := httpv1.NewBox(
		opts.RequestParsers.ExtendParser,
		opts.ResponseSender,
		useCase,
		metaOrderBy,
	)

	return controller, nil
}
