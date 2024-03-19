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

func initUnitFormElementEnvironment(ctx context.Context, opts factory.Options) (formElementOptions, error) {
	entityMetaUpdate, err := mrsql.NewEntityMetaUpdate(ctx, entity.FormElement{})

	if err != nil {
		return formElementOptions{}, err
	}

	storage := repository.NewFormElementPostgres(
		opts.PostgresAdapter,
		mrpostgres.NewSqlBuilderSelectCondition(
			mrpostgres.NewSqlBuilderWhere(),
		),
		mrpostgres.NewSqlBuilderUpdateWithMeta(
			entityMetaUpdate,
			mrpostgres.NewSqlBuilderSet(),
			nil,
		),
	)

	return formElementOptions{
		storage: storage,
	}, nil
}

func createUnitFormElement(ctx context.Context, opts moduleOptions) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitFormElement(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitFormElement(ctx context.Context, opts moduleOptions) (*http_v1.FormElement, error) {
	useCase := usecase.NewFormElement(
		opts.formElement.storage,
		opts.submitForm.storage,
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
