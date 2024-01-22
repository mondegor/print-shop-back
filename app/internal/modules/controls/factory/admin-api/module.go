package factory

import (
	module "print-shop-back/internal/modules/controls"
	"print-shop-back/internal/modules/controls/factory"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"
)

func CreateModule(opts *factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	mrfactory.InfoCreateModule(opts.Logger, module.Name)
	mrfactory.InfoCreateUnit(opts.Logger, module.UnitElementTemplateName)

	if l, err := createUnitElementTemplate(opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(l, module.UnitElementTemplatePermission)...)
	}

	mrfactory.InfoCreateUnit(opts.Logger, module.UnitFormDataName)

	if l, err := createUnitFormData(opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(l, module.UnitFormDataPermission)...)
	}

	mrfactory.InfoCreateUnit(opts.Logger, module.UnitFormElementName)

	if l, err := createUnitFormElement(opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(l, module.UnitFormElementPermission)...)
	}

	return list, nil
}
