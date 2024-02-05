package factory_filestation

import (
	"context"
	"print-shop-back/internal"
	"print-shop-back/internal/modules/file-station/factory"
)

func NewModuleOptions(ctx context.Context, opts app.Options) (factory.Options, error) {
	fileAPI, err := opts.FileProviderPool.Provider(
		opts.Cfg.ModulesSettings.FileStation.ImageProxy.FileProvider,
	)

	if err != nil {
		return factory.Options{}, err
	}

	return factory.Options{
		UsecaseHelper:  opts.UsecaseHelper,
		RequestParser:  opts.RequestParsers.String,
		ResponseSender: opts.ResponseSender,

		UnitImageProxy: factory.UnitImageProxyOptions{
			FileAPI: fileAPI,
			BaseURL: opts.Cfg.ModulesSettings.FileStation.ImageProxy.BaseURL,
		},
	}, nil
}
