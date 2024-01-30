package factory

import (
	"context"
	http_v1 "print-shop-back/internal/modules/file-station/controller/http_v1/public-api"
	"print-shop-back/internal/modules/file-station/factory"
	usecase "print-shop-back/internal/modules/file-station/usecase/public-api"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresponse"
)

func createUnitImageProxy(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitImageProxy(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitImageProxy(ctx context.Context, opts factory.Options) (*http_v1.ImageProxy, error) {
	service := usecase.NewFileProviderAdapter(opts.UnitImageProxy.FileAPI, opts.UsecaseHelper)
	controller := http_v1.NewImageProxy(
		opts.RequestParser,
		mrresponse.NewFileSender(opts.ResponseSender),
		service,
		opts.UnitImageProxy.BaseURL,
	)

	return controller, nil
}
