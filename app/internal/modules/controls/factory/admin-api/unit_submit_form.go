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

func createUnitSubmitForm(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitSubmitForm(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitSubmitForm(ctx context.Context, opts factory.Options) (*http_v1.SubmitForm, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.SubmitForm{})

	if err != nil {
		return nil, err
	}

	entityMetaUpdate, err := mrsql.NewEntityMetaUpdate(ctx, entity.SubmitForm{})

	if err != nil {
		return nil, err
	}

	storage := repository.NewSubmitFormPostgres(
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
	useCase := usecase.NewSubmitForm(storage, opts.EventEmitter, opts.UsecaseHelper)
	controller := http_v1.NewSubmitForm(
		opts.RequestParser,
		opts.ResponseSender,
		useCase,
		metaOrderBy,
	)

	return controller, nil
}
