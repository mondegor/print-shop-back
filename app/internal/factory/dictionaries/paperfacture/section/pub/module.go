package pub

import (
	"context"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/module"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/paperfacture"
)

// CreateModule - создаются все компоненты модуля и возвращаются к нему контролеры.
func CreateModule(ctx context.Context, opts paperfacture.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	mrfactory.InfoCreateModule(ctx, module.Name)

	if l, err := createUnitPaperFacture(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.PrepareEachController(l, mrfactory.WithPermission(module.Permission))...)
	}

	return list, nil
}
