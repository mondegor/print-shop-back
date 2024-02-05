package factory_catalog

import (
	"context"
	"print-shop-back/internal"
	view_shared "print-shop-back/internal/modules/catalog/box/controller/http_v1/shared/view"
	"print-shop-back/internal/modules/catalog/box/factory"
)

func NewBoxModuleOptions(ctx context.Context, opts app.Options) (factory.Options, error) {
	boxDictionary, err := opts.Translator.Dictionary("catalog/boxes")

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

		UnitBox: factory.UnitBoxOptions{
			Dictionary: boxDictionary,
		},
	}, nil
}
