package adm

import (
	"github.com/mondegor/go-components/factory/mrordering"
	"github.com/mondegor/go-storage/mrpostgres/builder"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/module"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/controls/submitform"
)

func initUnitFormElementEnvironment(opts submitform.Options) (formElementOptions, error) {
	entityMeta, err := mrsql.ParseEntity(opts.Logger, entity.FormElement{})
	if err != nil {
		return formElementOptions{}, err
	}

	storage := repository.NewFormElementPostgres(
		opts.DBConnManager,
		builder.NewSQL(
			builder.WithSQLSetMetaEntity(entityMeta.MetaUpdate()),
		),
	)

	return formElementOptions{
		storage: storage,
	}, nil
}

func createUnitFormElement(opts moduleOptions) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitFormElement(opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitFormElement(opts moduleOptions) (*httpv1.FormElement, error) { //nolint:unparam
	useCase := usecase.NewFormElement(
		opts.formElement.storage,
		opts.submitForm.storage,
		opts.ElementTemplateAPI,
		mrordering.NewComponentMover(
			opts.DBConnManager,
			mrsql.DBTableInfo{
				Name:       module.DBTableNameSubmitFormElements,
				PrimaryKey: "form_id",
			},
			opts.EventEmitter,
		),
		opts.EventEmitter,
		opts.UsecaseErrorWrapper,
		opts.Logger,
	)
	controller := httpv1.NewFormElement(
		opts.RequestParsers.ModuleParser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
