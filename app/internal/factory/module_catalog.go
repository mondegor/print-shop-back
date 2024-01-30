package factory

import (
	"context"
	"print-shop-back/internal/modules"
	view_shared "print-shop-back/internal/modules/catalog/controller/http_v1/shared/view"
	"print-shop-back/internal/modules/catalog/factory"
)

func NewCatalogModuleOptions(ctx context.Context, opts modules.Options) (factory.Options, error) {
	boxDictionary, err := opts.Translator.Dictionary("catalog/box")

	if err != nil {
		return factory.Options{}, err
	}

	laminateDictionary, err := opts.Translator.Dictionary("catalog/laminate")

	if err != nil {
		return factory.Options{}, err
	}

	paperDictionary, err := opts.Translator.Dictionary("catalog/paper")

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
			opts.RequestParsers.SortPage,
			opts.RequestParsers.String,
			opts.RequestParsers.Validator,
		),
		ResponseSender: opts.ResponseSender,

		LaminateTypeAPI: opts.DictionariesLaminateTypeAPI,
		PaperColorAPI:   opts.DictionariesPaperColorAPI,
		PaperFactureAPI: opts.DictionariesPaperFactureAPI,

		UnitBox: factory.UnitBoxOptions{
			Dictionary: boxDictionary,
		},

		UnitLaminate: factory.UnitLaminateOptions{
			Dictionary: laminateDictionary,
		},

		UnitPaper: factory.UnitPaperOptions{
			Dictionary: paperDictionary,
		},
	}, nil
}
