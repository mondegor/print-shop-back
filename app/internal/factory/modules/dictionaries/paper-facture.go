package factory_dictionaries

import (
	"context"
	"print-shop-back/internal"
	view_shared "print-shop-back/internal/modules/dictionaries/paper-facture/controller/http_v1/shared/view"
	"print-shop-back/internal/modules/dictionaries/paper-facture/factory"
	factory_api "print-shop-back/internal/modules/dictionaries/paper-facture/factory/api"
	usecase_api "print-shop-back/internal/modules/dictionaries/paper-facture/usecase/api"

	"github.com/mondegor/go-webcore/mrlog"
)

func NewPaperFactureModuleOptions(ctx context.Context, opts app.Options) (factory.Options, error) {
	paperFactureDictionary, err := opts.Translator.Dictionary("dictionaries/paper-factures")

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

		UnitPaperFacture: factory.UnitPaperFactureOptions{
			Dictionary: paperFactureDictionary,
		},

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}

func NewPaperFactureAPI(ctx context.Context, opts app.Options) (*usecase_api.PaperFacture, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init dictionaries paper facture API")

	return factory_api.NewPaperFacture(opts.PostgresAdapter, opts.UsecaseHelper), nil
}
