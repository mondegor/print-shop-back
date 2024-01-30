package factory

import (
	"context"
	module "print-shop-back/internal/modules/catalog"
	"print-shop-back/internal/modules/catalog/factory"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"
)

func CreateModule(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	mrfactory.InfoCreateModule(ctx, module.Name)
	mrfactory.InfoCreateUnit(ctx, module.UnitBoxName)

	if l, err := createUnitBox(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(ctx, l, module.UnitBoxPermission)...)
	}

	mrfactory.InfoCreateUnit(ctx, module.UnitLaminateName)

	if l, err := createUnitLaminate(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(ctx, l, module.UnitLaminatePermission)...)
	}

	mrfactory.InfoCreateUnit(ctx, module.UnitPaperName)

	if l, err := createUnitPaper(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(ctx, l, module.UnitPaperPermission)...)
	}

	return list, nil
}
