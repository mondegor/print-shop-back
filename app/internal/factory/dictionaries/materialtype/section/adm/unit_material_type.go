package adm

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/adm/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/adm/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/materialtype"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
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
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.MaterialType{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewMaterialTypePostgres(
		opts.DBConnManager,
		mrpostgres.NewSQLBuilderSelect(
			mrpostgres.NewSQLBuilderWhere(),
			mrpostgres.NewSQLBuilderOrderBy(ctx, metaOrderBy.DefaultSort()),
			mrpostgres.NewSQLBuilderLimit(opts.PageSizeMax),
		),
	)
	useCase := usecase.NewMaterialType(storage, opts.EventEmitter, opts.UsecaseHelper)
	controller := httpv1.NewMaterialType(
		opts.RequestParsers.ExtendParser,
		opts.ResponseSender,
		useCase,
		metaOrderBy,
	)

	return controller, nil
}
