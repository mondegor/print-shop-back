package factory

import (
	"context"
	"print-shop-back/internal/modules"
	"print-shop-back/internal/modules/file-station/factory"
)

func NewFileStationModuleOptions(ctx context.Context, opts modules.Options) (factory.Options, error) {
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
