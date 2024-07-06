package adm

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/controls/submitform"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
)

func initUnitFormElementEnvironment(ctx context.Context, opts submitform.Options) (formElementOptions, error) {
	entityMetaUpdate, err := mrsql.NewEntityMetaUpdate(ctx, entity.FormElement{})
	if err != nil {
		return formElementOptions{}, err
	}

	storage := repository.NewFormElementPostgres(
		opts.DBConnManager,
		mrpostgres.NewSQLBuilderCondition(
			mrpostgres.NewSQLBuilderWhere(),
		),
		mrpostgres.NewSQLBuilderUpdateWithMeta(
			entityMetaUpdate,
			mrpostgres.NewSQLBuilderSet(),
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

func newUnitFormElement(_ context.Context, opts moduleOptions) (*httpv1.FormElement, error) { //nolint:unparam
	useCase := usecase.NewFormElement(
		opts.formElement.storage,
		opts.submitForm.storage,
		opts.ElementTemplateAPI,
		opts.OrdererAPI,
		opts.EventEmitter,
		opts.UsecaseHelper,
	)
	controller := httpv1.NewFormElement(
		opts.RequestParsers.ModuleParser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
