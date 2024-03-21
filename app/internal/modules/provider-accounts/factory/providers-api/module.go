package factory

import (
	"context"
	module "print-shop-back/internal/modules/provider-accounts"
	"print-shop-back/internal/modules/provider-accounts/factory"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"
)

func CreateModule(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	mrfactory.InfoCreateModule(ctx, module.Name)

	if l, err := createUnitCompanyPage(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(ctx, l, module.UnitCompanyPagePermission)...)
	}

	return list, nil
}
