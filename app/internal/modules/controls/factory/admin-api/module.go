package factory

import (
	"context"
	module "print-shop-back/internal/modules/controls"
	"print-shop-back/internal/modules/controls/factory"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"
)

func CreateModule(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	mrfactory.InfoCreateModule(ctx, module.Name)
	mrfactory.InfoCreateUnit(ctx, module.UnitElementTemplateName)

	if l, err := createUnitElementTemplate(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(ctx, l, module.UnitElementTemplatePermission)...)
	}

	mrfactory.InfoCreateUnit(ctx, module.UnitSubmitFormName)

	if l, err := createUnitSubmitForm(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(ctx, l, module.UnitSubmitFormPermission)...)
	}

	mrfactory.InfoCreateUnit(ctx, module.UnitFormElementName)

	if l, err := createUnitFormElement(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(ctx, l, module.UnitFormElementPermission)...)
	}

	return list, nil
}
