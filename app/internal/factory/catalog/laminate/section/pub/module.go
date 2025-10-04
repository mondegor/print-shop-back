package pub

import (
	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/catalog/laminate/module"
	"github.com/mondegor/print-shop-back/internal/factory/catalog/laminate"
)

// CreateModule - создаются все компоненты модуля и возвращаются к нему контролеры.
func CreateModule(opts laminate.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	mrinit.InfoCreateModule(opts.Logger, module.Name)

	if l, err := createUnitLaminate(opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrinit.PrepareEachController(l, mrinit.WithPermission(module.Permission))...)
	}

	return list, nil
}
