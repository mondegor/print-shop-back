package pub

import (
	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/module"
	"github.com/mondegor/print-shop-back/internal/factory/calculations/algo"
)

// CreateModule - создаются все компоненты модуля и возвращаются к нему контролеры.
func CreateModule(opts algo.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	mrinit.InfoCreateModule(opts.Logger, module.Name)

	if l, err := createUnitBoxPackInBox(opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrinit.PrepareEachController(l, mrinit.WithPermission(module.Permission))...)
	}

	if l, err := createUnitSheetPackInStack(opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrinit.PrepareEachController(l, mrinit.WithPermission(module.Permission))...)
	}

	if l, err := createUnitSheetCutting(opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrinit.PrepareEachController(l, mrinit.WithPermission(module.Permission))...)
	}

	if l, err := createUnitSheetImposition(opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrinit.PrepareEachController(l, mrinit.WithPermission(module.Permission))...)
	}

	if l, err := createUnitSheetInsideOutside(opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrinit.PrepareEachController(l, mrinit.WithPermission(module.Permission))...)
	}

	return list, nil
}
