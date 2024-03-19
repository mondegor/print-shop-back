package factory

import (
	"context"
	http_v1 "print-shop-back/internal/modules/controls/element-template/controller/http_v1/admin-api"
	entity "print-shop-back/internal/modules/controls/element-template/entity/admin-api"
	"print-shop-back/internal/modules/controls/element-template/factory"
	repository "print-shop-back/internal/modules/controls/element-template/infrastructure/repository/admin-api"
	usecase "print-shop-back/internal/modules/controls/element-template/usecase/admin-api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresponse"
)

func createUnitElementTemplate(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitElementTemplate(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitElementTemplate(ctx context.Context, opts factory.Options) (*http_v1.ElementTemplate, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.ElementTemplate{})

	if err != nil {
		return nil, err
	}

	entityMetaUpdate, err := mrsql.NewEntityMetaUpdate(ctx, entity.ElementTemplate{})

	if err != nil {
		return nil, err
	}

	storage := repository.NewElementTemplatePostgres(
		opts.PostgresAdapter,
		mrpostgres.NewSqlBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			mrpostgres.NewSqlBuilderOrderBy(ctx, metaOrderBy.DefaultSort()),
			mrpostgres.NewSqlBuilderPager(opts.PageSizeMax),
		),
		mrpostgres.NewSqlBuilderUpdateWithMeta(
			entityMetaUpdate,
			mrpostgres.NewSqlBuilderSet(),
			nil,
		),
	)
	useCase := usecase.NewElementTemplate(storage, opts.EventEmitter, opts.UsecaseHelper)
	controller := http_v1.NewElementTemplate(
		opts.RequestParser,
		mrresponse.NewFileSender(opts.ResponseSender),
		useCase,
		metaOrderBy,
	)

	return controller, nil
}
