package factory

import (
	"context"
	http_v1 "print-shop-back/internal/modules/controls/controller/http_v1/admin-api"
	entity "print-shop-back/internal/modules/controls/entity/admin-api"
	"print-shop-back/internal/modules/controls/factory"
	repository "print-shop-back/internal/modules/controls/infrastructure/repository/admin-api"
	usecase "print-shop-back/internal/modules/controls/usecase/admin-api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitFormElement(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitFormElement(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitFormElement(ctx context.Context, opts factory.Options) (*http_v1.FormElement, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.FormElement{})

	if err != nil {
		return nil, err
	}

	entityMetaUpdate, err := mrsql.NewEntityMetaUpdate(ctx, entity.FormElement{})

	if err != nil {
		return nil, err
	}

	storage := repository.NewFormElementPostgres(
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
	useCase := usecase.NewFormElement(storage, opts.ElementTemplateAPI, opts.OrdererAPI, opts.EventEmitter, opts.UsecaseHelper)
	controller := http_v1.NewFormElement(
		opts.RequestParser,
		opts.ResponseSender,
		useCase,
		metaOrderBy,
	)

	return controller, nil
}
