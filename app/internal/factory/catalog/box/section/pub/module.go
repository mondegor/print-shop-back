package pub

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/catalog/box/module"
	"github.com/mondegor/print-shop-back/internal/factory/catalog/box"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"
)

// CreateModule - создаются все компоненты модуля и возвращаются к нему контролеры.
func CreateModule(ctx context.Context, opts box.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	mrfactory.InfoCreateModule(ctx, module.Name)

	if l, err := createUnitBox(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.PrepareEachController(l, mrfactory.WithPermission(module.Permission))...)
	}

	return list, nil
}
