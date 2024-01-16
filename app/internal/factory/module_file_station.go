package factory

import (
	"print-shop-back/internal/modules"
	"print-shop-back/internal/modules/file-station/factory"
)

func NewFileStationOptions(opts *modules.Options) (*factory.Options, error) {
	fileAPI, err := opts.FileProviderPool.Provider(
		opts.Cfg.ModulesSettings.FileStation.ImageProxy.FileProvider,
	)

	if err != nil {
		return nil, err
	}

	return &factory.Options{
		Logger:        opts.Logger,
		ServiceHelper: opts.ServiceHelper,

		UnitImageProxy: &factory.UnitImageProxyOptions{
			FileAPI: fileAPI,
			BaseURL: opts.Cfg.ModulesSettings.FileStation.ImageProxy.BaseURL,
		},
	}, nil
}
