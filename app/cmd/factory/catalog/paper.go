package catalog

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore/mrinit"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/catalog/paper/module"
	"github.com/mondegor/print-shop-back/internal/factory/catalog/paper"
)

// NewPaperModuleOptions - создаёт объект paper.Options.
func NewPaperModuleOptions(_ context.Context, opts app.Options) (paper.Options, error) {
	paperDictionary, err := opts.Translator.Dictionary("catalog/papers")
	if err != nil {
		return paper.Options{}, err
	}

	return paper.Options{
		EventEmitter:  opts.EventEmitter,
		UsecaseHelper: opts.UsecaseErrorWrapper,
		DBConnManager: opts.PostgresConnManager,
		RequestParsers: paper.RequestParsers{
			Parser:       opts.RequestParsers.Parser,
			ExtendParser: opts.RequestParsers.ExtendParser,
		},
		ResponseSender: opts.ResponseSenders.Sender,

		MaterialTypeAPI: opts.DictionariesMaterialTypeAPI,
		PaperColorAPI:   opts.DictionariesPaperColorAPI,
		PaperFactureAPI: opts.DictionariesPaperFactureAPI,

		UnitPaper: paper.UnitPaperOptions{
			Dictionary: paperDictionary,
		},

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}

// RegisterPaperErrors - comment func.
func RegisterPaperErrors(em *mrinit.ErrorManager) {
	em.RegisterList(mrinit.WrapProtoList(module.Errors()))
}
