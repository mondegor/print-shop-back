package factory

import (
	module "print-shop-back/internal/modules/catalog"
	"print-shop-back/internal/modules/catalog/factory"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"
)

func CreateModule(opts *factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	mrfactory.InfoCreateModule(opts.Logger, module.Name)
	mrfactory.InfoCreateUnit(opts.Logger, module.UnitBoxName)

	if l, err := createUnitBox(opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(l, module.UnitBoxPermission)...)
	}

	mrfactory.InfoCreateUnit(opts.Logger, module.UnitLaminateName)

	if l, err := createUnitLaminate(opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(l, module.UnitLaminatePermission)...)
	}

	mrfactory.InfoCreateUnit(opts.Logger, module.UnitPaperName)

	if l, err := createUnitPaper(opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(l, module.UnitPaperPermission)...)
	}

	return list, nil
}
