package filestation

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/factory/filestation"
)

// NewModuleOptions - создаёт объект filestation.Options.
func NewModuleOptions(_ context.Context, opts app.Options) (filestation.Options, error) {
	fileAPI, err := opts.FileProviderPool.Provider(
		opts.Cfg.ModulesSettings.FileStation.ImageProxy.FileProvider,
	)
	if err != nil {
		return filestation.Options{}, err
	}

	return filestation.Options{
		UsecaseHelper:  opts.UsecaseErrorWrapper,
		RequestParser:  opts.RequestParsers.String,
		ResponseSender: opts.ResponseSenders.FileSender,

		UnitImageProxy: filestation.UnitImageProxyOptions{
			FileAPI: fileAPI,
			BaseURL: opts.Cfg.ModulesSettings.FileStation.ImageProxy.BaseURL,
		},
	}, nil
}
