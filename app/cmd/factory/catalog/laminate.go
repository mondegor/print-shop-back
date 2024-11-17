package catalog

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore/mrapp"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/factory/catalog/laminate"
)

// NewLaminateModuleOptions - создаёт объект laminate.Options.
func NewLaminateModuleOptions(_ context.Context, opts app.Options) (laminate.Options, error) {
	laminateDictionary, err := opts.Translator.Dictionary("catalog/laminates")
	if err != nil {
		return laminate.Options{}, err
	}

	return laminate.Options{
		EventEmitter:        opts.EventEmitter,
		UseCaseErrorWrapper: mrapp.NewUseCaseErrorWrapper(),
		DBConnManager:       opts.PostgresConnManager,
		RequestParsers: laminate.RequestParsers{
			Parser:       opts.RequestParsers.Parser,
			ExtendParser: opts.RequestParsers.ExtendParser,
		},
		ResponseSender: opts.ResponseSenders.Sender,

		MaterialTypeAPI: opts.DictionariesMaterialTypeAPI,

		UnitLaminate: laminate.UnitLaminateOptions{
			Dictionary: laminateDictionary,
		},

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}
