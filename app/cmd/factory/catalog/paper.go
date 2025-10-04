package catalog

import (
	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/factory/catalog/paper"
)

// NewPaperModuleOptions - создаёт объект paper.Options.
func NewPaperModuleOptions(opts app.Options) paper.Options {
	return paper.Options{
		Logger:              opts.Logger,
		EventEmitter:        opts.EventEmitter,
		UsecaseErrorWrapper: opts.UsecaseErrorWrapper,
		DBConnManager:       opts.PostgresConnManager,
		RequestParsers: paper.RequestParsers{
			Parser:       opts.RequestParsers.Parser,
			ExtendParser: opts.RequestParsers.ExtendParser,
		},
		ResponseSender: opts.ResponseSenders.Sender,

		MaterialTypeAPI: opts.DictionariesMaterialTypeAPI,
		PaperColorAPI:   opts.DictionariesPaperColorAPI,
		PaperFactureAPI: opts.DictionariesPaperFactureAPI,

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}
}
