package factory

import (
	module "print-shop-back/internal/modules/dictionaries"
	"print-shop-back/internal/modules/dictionaries/factory"

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
	opts.Logger.Info("Init unit %s in %s section", module.UnitLaminateTypeName, section.Caption())

	if err := newUnitLaminateType(c, opts, section); err != nil {
		return err
	}

	opts.Logger.Info("Init unit %s in %s section", module.UnitPaperColorName, section.Caption())

	if err := newUnitPaperColor(c, opts, section); err != nil {
		return err
	}

	opts.Logger.Info("Init unit %s in %s section", module.UnitPaperFactureName, section.Caption())

	if err := newUnitPaperFacture(c, opts, section); err != nil {
		return err
	}

	opts.Logger.Info("Init unit %s in %s section", module.UnitPrintFormatName, section.Caption())

	if err := newUnitPrintFormat(c, opts, section); err != nil {
		return err
	}

	return nil
}
