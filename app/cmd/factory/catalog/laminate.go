package catalog

import (
	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/factory/catalog/laminate"
)

// NewLaminateModuleOptions - создаёт объект laminate.Options.
func NewLaminateModuleOptions(opts app.Options) laminate.Options {
	return laminate.Options{
		Logger:              opts.Logger,
		EventEmitter:        opts.EventEmitter,
		UsecaseErrorWrapper: opts.UsecaseErrorWrapper,
		DBConnManager:       opts.PostgresConnManager,
		RequestParsers: laminate.RequestParsers{
			Parser:       opts.RequestParsers.Parser,
			ExtendParser: opts.RequestParsers.ExtendParser,
		},
		ResponseSender: opts.ResponseSenders.Sender,

		MaterialTypeAPI: opts.DictionariesMaterialTypeAPI,

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}
}
