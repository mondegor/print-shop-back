package factory

import (
	"context"
	module "print-shop-back/internal/modules/controls/submit-form"
	"print-shop-back/internal/modules/controls/submit-form/factory"
	repository "print-shop-back/internal/modules/controls/submit-form/infrastructure/repository/admin-api"

	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"
)

type (
	moduleOptions struct {
		factory.Options

		submitForm  submitFormOptions
		formElement formElementOptions
	}

	submitFormOptions struct {
		metaOrderBy *mrsql.EntityMetaOrderBy
		storage     *repository.SubmitFormPostgres
	}

	formElementOptions struct {
		storage *repository.FormElementPostgres
	}
)

func CreateModule(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	mrfactory.InfoCreateModule(ctx, module.Name)

	unitSubmitFormOptions, err := initUnitSubmitFormEnvironment(ctx, opts)

	if err != nil {
		return nil, err
	}

	unitFormElementOptions, err := initUnitFormElementEnvironment(ctx, opts)

	if err != nil {
		return nil, err
	}

	return createModule(
		ctx,
		moduleOptions{
			Options:     opts,
			submitForm:  unitSubmitFormOptions,
			formElement: unitFormElementOptions,
		},
	)
}

func createModule(ctx context.Context, opts moduleOptions) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if l, err := createUnitSubmitForm(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(ctx, l, module.UnitSubmitFormPermission)...)
	}

	if l, err := createUnitFormElement(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(ctx, l, module.UnitFormElementPermission)...)
	}

	return list, nil
}
