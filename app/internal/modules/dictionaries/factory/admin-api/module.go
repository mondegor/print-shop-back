package factory

import (
	module "print-shop-back/internal/modules/dictionaries"
	"print-shop-back/internal/modules/dictionaries/factory"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"
)

func CreateModule(opts *factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	mrfactory.InfoCreateModule(opts.Logger, module.Name)
	mrfactory.InfoCreateUnit(opts.Logger, module.UnitLaminateTypeName)

	if l, err := createUnitLaminateType(opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(l, module.UnitLaminateTypePermission)...)
	}

	mrfactory.InfoCreateUnit(opts.Logger, module.UnitPaperColorName)

	if l, err := createUnitPaperColor(opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(l, module.UnitPaperColorPermission)...)
	}

	mrfactory.InfoCreateUnit(opts.Logger, module.UnitPaperFactureName)

	if l, err := createUnitPaperFacture(opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(l, module.UnitPaperFacturePermission)...)
	}

	mrfactory.InfoCreateUnit(opts.Logger, module.UnitPrintFormatName)

	if l, err := createUnitPrintFormat(opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(l, module.UnitPrintFormatPermission)...)
	}

	return list, nil
}
