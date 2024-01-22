package factory

import (
	module "print-shop-back/internal/modules/provider-accounts"
	"print-shop-back/internal/modules/provider-accounts/factory"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"
)

func CreateModule(opts *factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	mrfactory.InfoCreateModule(opts.Logger, module.Name)
	mrfactory.InfoCreateUnit(opts.Logger, module.UnitCompanyPageName)

	if l, err := createUnitCompanyPage(opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(l, module.UnitCompanyPagePermission)...)
	}

	return list, nil
}
