package adm

import (
	"context"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/section/adm/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/section/adm/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/controls/elementtemplate"
)

func createUnitElementTemplate(ctx context.Context, opts elementtemplate.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitElementTemplate(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitElementTemplate(ctx context.Context, opts elementtemplate.Options) (*httpv1.ElementTemplate, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.ElementTemplate{})
	if err != nil {
		return nil, err
	}

	entityMetaUpdate, err := mrsql.NewEntityMetaUpdate(ctx, entity.ElementTemplate{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewElementTemplatePostgres(
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
	useCase := usecase.NewElementTemplate(storage, opts.EventEmitter, opts.UseCaseHelper)
	controller := httpv1.NewElementTemplate(
		opts.RequestParsers.ModuleParser,
		opts.ResponseSender,
		useCase,
		metaOrderBy,
	)

	return controller, nil
}
