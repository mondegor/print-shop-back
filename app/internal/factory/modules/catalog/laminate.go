package factory_catalog

import (
	"context"
	"print-shop-back/internal"
	view_shared "print-shop-back/internal/modules/catalog/laminate/controller/http_v1/shared/view"
	"print-shop-back/internal/modules/catalog/laminate/factory"
)

func NewLaminateModuleOptions(ctx context.Context, opts app.Options) (factory.Options, error) {
	laminateDictionary, err := opts.Translator.Dictionary("catalog/laminates")

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

		LaminateTypeAPI: opts.DictionariesLaminateTypeAPI,

		UnitLaminate: factory.UnitLaminateOptions{
			Dictionary: laminateDictionary,
		},

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}
