package factory

import (
	module "print-shop-back/internal/modules/provider-accounts"
	"print-shop-back/internal/modules/provider-accounts/factory"

	"github.com/mondegor/go-webcore/mrcore"
)

func NewModule(opts *factory.Options, section mrcore.ClientSection) ([]mrcore.HttpController, error) {
	opts.Logger.Info("Init module %s in section %s", module.Name, section.Caption())

	var c []mrcore.HttpController

	if err := newModule(&c, opts, section); err != nil {
		return nil, err
	}

	return c, nil
}

func newModule(c *[]mrcore.HttpController, opts *factory.Options, section mrcore.ClientSection) error {
	opts.Logger.Info("Init unit %s in %s section", module.UnitCompanyPageName, section.Caption())

	if err := newUnitCompanyPage(c, opts, section); err != nil {
		return err
	}

	return nil
}
