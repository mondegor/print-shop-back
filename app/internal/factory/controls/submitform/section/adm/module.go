package adm

import (
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrcore/mrinit"
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
func CreateModule(opts submitform.Options) ([]mrserver.HttpController, error) {
	mrinit.InfoCreateModule(opts.Logger, module.Name)

	unitSubmitFormOptions, err := initUnitSubmitFormEnvironment(opts)
	if err != nil {
		return nil, err
	}

	unitFormElementOptions, err := initUnitFormElementEnvironment(opts)
	if err != nil {
		return nil, err
	}

	unitFormVersionOptions, err := initUnitSubmitFormVersionEnvironment(opts)
	if err != nil {
		return nil, err
	}

	return createModule(
		moduleOptions{
			Options:     opts,
			submitForm:  unitSubmitFormOptions,
			formElement: unitFormElementOptions,
			formVersion: unitFormVersionOptions,
		},
	)
}

func createModule(opts moduleOptions) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if l, err := createUnitSubmitForm(opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrinit.PrepareEachController(l, mrinit.WithPermission(module.UnitSubmitFormPermission))...)
	}

	if l, err := createUnitFormElement(opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrinit.PrepareEachController(l, mrinit.WithPermission(module.UnitFormElementPermission))...)
	}

	return list, nil
}
