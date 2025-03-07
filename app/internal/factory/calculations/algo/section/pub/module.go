package pub

import (
	"context"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/module"
	"github.com/mondegor/print-shop-back/internal/factory/calculations/algo"
)

// CreateModule - создаются все компоненты модуля и возвращаются к нему контролеры.
func CreateModule(ctx context.Context, opts algo.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	mrfactory.InfoCreateModule(ctx, module.Name)

	if l, err := createUnitBoxPackInBox(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.PrepareEachController(l, mrfactory.WithPermission(module.Permission))...)
	}

	if l, err := createUnitSheetPackInStack(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.PrepareEachController(l, mrfactory.WithPermission(module.Permission))...)
	}

	if l, err := createUnitSheetCutting(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.PrepareEachController(l, mrfactory.WithPermission(module.Permission))...)
	}

	if l, err := createUnitSheetImposition(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.PrepareEachController(l, mrfactory.WithPermission(module.Permission))...)
	}

	if l, err := createUnitSheetInsideOutside(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.PrepareEachController(l, mrfactory.WithPermission(module.Permission))...)
	}

	return list, nil
}
