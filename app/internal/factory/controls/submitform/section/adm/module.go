package adm

import (
	"context"

	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/module"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/factory/controls/submitform"
)

type (
	moduleOptions struct {
		submitform.Options

		submitForm  submitFormOptions
		formElement formElementOptions
		formVersion formVersionOptions
	}

	submitFormOptions struct {
		metaOrderBy *mrsql.EntityMetaOrderBy
		storage     *repository.SubmitFormPostgres
	}

	formElementOptions struct {
		storage *repository.FormElementPostgres
	}

	formVersionOptions struct {
		storage *repository.FormVersionPostgres
	}
)

// CreateModule - создаются все компоненты модуля и возвращаются к нему контролеры.
func CreateModule(ctx context.Context, opts submitform.Options) ([]mrserver.HttpController, error) {
	mrfactory.InfoCreateModule(ctx, module.Name)

	unitSubmitFormOptions, err := initUnitSubmitFormEnvironment(ctx, opts)
	if err != nil {
		return nil, err
	}

	unitFormElementOptions, err := initUnitFormElementEnvironment(ctx, opts)
	if err != nil {
		return nil, err
	}

	unitFormVersionOptions, err := initUnitSubmitFormVersionEnvironment(ctx, opts)
	if err != nil {
		return nil, err
	}

	return createModule(
		ctx,
		moduleOptions{
			Options:     opts,
			submitForm:  unitSubmitFormOptions,
			formElement: unitFormElementOptions,
			formVersion: unitFormVersionOptions,
		},
	)
}

func createModule(ctx context.Context, opts moduleOptions) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if l, err := createUnitSubmitForm(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.PrepareEachController(l, mrfactory.WithPermission(module.UnitSubmitFormPermission))...)
	}

	if l, err := createUnitFormElement(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.PrepareEachController(l, mrfactory.WithPermission(module.UnitFormElementPermission))...)
	}

	return list, nil
}
