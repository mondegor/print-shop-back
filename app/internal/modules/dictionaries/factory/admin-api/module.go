package factory

import (
	"context"
	module "print-shop-back/internal/modules/dictionaries"
	"print-shop-back/internal/modules/dictionaries/factory"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"
)

func CreateModule(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	mrfactory.InfoCreateModule(ctx, module.Name)
	mrfactory.InfoCreateUnit(ctx, module.UnitLaminateTypeName)

	if l, err := createUnitLaminateType(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(ctx, l, module.UnitLaminateTypePermission)...)
	}

	mrfactory.InfoCreateUnit(ctx, module.UnitPaperColorName)

	if l, err := createUnitPaperColor(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(ctx, l, module.UnitPaperColorPermission)...)
	}

	mrfactory.InfoCreateUnit(ctx, module.UnitPaperFactureName)

	if l, err := createUnitPaperFacture(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(ctx, l, module.UnitPaperFacturePermission)...)
	}

	mrfactory.InfoCreateUnit(ctx, module.UnitPrintFormatName)

	if l, err := createUnitPrintFormat(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(ctx, l, module.UnitPrintFormatPermission)...)
	}

	return list, nil
}
