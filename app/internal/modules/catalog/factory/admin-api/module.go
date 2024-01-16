package factory

import (
	module "print-shop-back/internal/modules/catalog"
	"print-shop-back/internal/modules/catalog/factory"

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
	opts.Logger.Info("Init unit %s in %s section", module.UnitBoxName, section.Caption())

	if err := newUnitBox(c, opts, section); err != nil {
		return err
	}

	opts.Logger.Info("Init unit %s in %s section", module.UnitLaminateName, section.Caption())

	if err := newUnitLaminate(c, opts, section); err != nil {
		return err
	}

	opts.Logger.Info("Init unit %s in %s section", module.UnitPaperName, section.Caption())

	if err := newUnitPaper(c, opts, section); err != nil {
		return err
	}

	return nil
}
