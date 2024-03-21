package factory_dictionaries

import (
	"context"
	"print-shop-back/internal"
	view_shared "print-shop-back/internal/modules/dictionaries/print-format/controller/http_v1/shared/view"
	"print-shop-back/internal/modules/dictionaries/print-format/factory"
	factory_api "print-shop-back/internal/modules/dictionaries/print-format/factory/api"
	usecase_api "print-shop-back/internal/modules/dictionaries/print-format/usecase/api"

	"github.com/mondegor/go-webcore/mrlog"
)

func NewPrintFormatModuleOptions(ctx context.Context, opts app.Options) (factory.Options, error) {
	printFormatDictionary, err := opts.Translator.Dictionary("dictionaries/print-formats")

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

		UnitPrintFormat: factory.UnitPrintFormatOptions{
			Dictionary: printFormatDictionary,
		},

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}

func NewPrintFormatAPI(ctx context.Context, opts app.Options) (*usecase_api.PrintFormat, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init dictionaries print format API")

	return factory_api.NewPrintFormat(opts.PostgresAdapter, opts.UsecaseHelper), nil
}
