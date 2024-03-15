package factory_catalog

import (
	"context"
	"print-shop-back/internal"
	view_shared "print-shop-back/internal/modules/catalog/paper/controller/http_v1/shared/view"
	"print-shop-back/internal/modules/catalog/paper/factory"
)

func NewPaperModuleOptions(ctx context.Context, opts app.Options) (factory.Options, error) {
	paperDictionary, err := opts.Translator.Dictionary("catalog/papers")

	if err != nil {
		return factory.Options{}, err
	}

	return factory.Options{
		EventEmitter:    opts.EventEmitter,
		UsecaseHelper:   opts.UsecaseHelper,
		PostgresAdapter: opts.PostgresAdapter,
		RequestParser: view_shared.NewParser(
			opts.RequestParsers.Int64,
			opts.RequestParsers.ItemStatus,
			opts.RequestParsers.KeyInt32,
			opts.RequestParsers.ListSorter,
			opts.RequestParsers.ListPager,
			opts.RequestParsers.String,
			opts.RequestParsers.Validator,
		),
		ResponseSender: opts.ResponseSender,

		PaperColorAPI:   opts.DictionariesPaperColorAPI,
		PaperFactureAPI: opts.DictionariesPaperFactureAPI,

		UnitPaper: factory.UnitPaperOptions{
			Dictionary: paperDictionary,
		},

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}
