package adm

import (
	"context"

	"github.com/mondegor/go-storage/mrpostgres/builder"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/adm/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/adm/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/materialtype"
)

func createUnitMaterialType(ctx context.Context, opts materialtype.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitMaterialType(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitMaterialType(ctx context.Context, opts materialtype.Options) (*httpv1.MaterialType, error) {
	entityMeta, err := mrsql.ParseEntity(mrlog.Ctx(ctx), entity.MaterialType{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewMaterialTypePostgres(
		opts.DBConnManager,
		builder.NewSQL(
			builder.WithSQLOrderByDefaultSort(entityMeta.MetaOrderBy().DefaultSort()),
			builder.WithSQLLimitMaxSize(opts.PageSizeMax),
		),
	)
	useCase := usecase.NewMaterialType(storage, opts.EventEmitter, opts.UseCaseErrorWrapper)
	controller := httpv1.NewMaterialType(
		opts.RequestParsers.ExtendParser,
		opts.ResponseSender,
		useCase,
		entityMeta.MetaOrderBy(),
	)

	return controller, nil
}
