package filestation

import (
	"context"

	"github.com/mondegor/go-webcore/mrpath/placeholderpath"

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

	basePath, err := placeholderpath.New(
		opts.Cfg.ModulesSettings.FileStation.ImageProxy.BasePath,
		placeholderpath.Placeholder,
	)
	if err != nil {
		return filestation.Options{}, err
	}

	return filestation.Options{
		UsecaseHelper:  opts.UsecaseErrorWrapper,
		RequestParser:  opts.RequestParsers.String,
		ResponseSender: opts.ResponseSenders.FileSender,

		UnitImageProxy: filestation.UnitImageProxyOptions{
			FileAPI:  fileAPI,
			BasePath: basePath,
		},
	}, nil
}
