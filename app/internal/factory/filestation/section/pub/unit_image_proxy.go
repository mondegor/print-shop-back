package pub

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/factory/filestation"
	"github.com/mondegor/print-shop-back/internal/filestation/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/filestation/section/pub/usecase"

	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitImageProxy(ctx context.Context, opts filestation.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitImageProxy(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitImageProxy(_ context.Context, opts filestation.Options) (*httpv1.ImageProxy, error) { //nolint:unparam
	useCase := usecase.NewFileProviderAdapter(opts.UnitImageProxy.FileAPI, opts.UsecaseHelper)
	controller := httpv1.NewImageProxy(
		opts.RequestParser,
		opts.ResponseSender,
		useCase,
		opts.UnitImageProxy.BaseURL,
	)

	return controller, nil
}
