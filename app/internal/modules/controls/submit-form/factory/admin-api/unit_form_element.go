package factory

import (
	"context"
	http_v1 "print-shop-back/internal/modules/controls/submit-form/controller/http_v1/admin-api"
	entity "print-shop-back/internal/modules/controls/submit-form/entity/admin-api"
	"print-shop-back/internal/modules/controls/submit-form/factory"
	repository "print-shop-back/internal/modules/controls/submit-form/infrastructure/repository/admin-api"
	usecase "print-shop-back/internal/modules/controls/submit-form/usecase/admin-api"

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
	_, storageForm, err := newUnitSubmitFormStorage(ctx, opts)

	if err != nil {
		return nil, err
	}

	storage, err := newUnitFormElementStorage(ctx, opts)

	if err != nil {
		return nil, err
	}

	useCase := usecase.NewFormElement(
		storage,
		storageForm,
		opts.ElementTemplateAPI,
		opts.OrdererAPI,
		opts.EventEmitter,
		opts.UsecaseHelper,
	)
	controller := http_v1.NewFormElement(
		opts.RequestParser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}

func newUnitFormElementStorage(ctx context.Context, opts factory.Options) (*repository.FormElementPostgres, error) {
	entityMetaUpdate, err := mrsql.NewEntityMetaUpdate(ctx, entity.FormElement{})

	if err != nil {
		return nil, err
	}

	return repository.NewFormElementPostgres(
		opts.PostgresAdapter,
		mrpostgres.NewSqlBuilderSelectCondition(
			mrpostgres.NewSqlBuilderWhere(),
		),
		mrpostgres.NewSqlBuilderUpdateWithMeta(
			entityMetaUpdate,
			mrpostgres.NewSqlBuilderSet(),
			nil,
		),
	), nil
}
