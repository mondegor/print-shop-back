package factory

import (
	"context"
	"print-shop-back/internal/modules"
	view_shared "print-shop-back/internal/modules/dictionaries/controller/http_v1/shared/view"
	"print-shop-back/internal/modules/dictionaries/factory"
	factory_api "print-shop-back/internal/modules/dictionaries/factory/api"
	usecase_api "print-shop-back/internal/modules/dictionaries/usecase/api"

	"github.com/mondegor/go-webcore/mrlog"
)

func NewDictionariesModuleOptions(ctx context.Context, opts modules.Options) (factory.Options, error) {
	laminateTypeDictionary, err := opts.Translator.Dictionary("dictionaries/laminate-type")

	if err != nil {
		return factory.Options{}, err
	}

	paperColorDictionary, err := opts.Translator.Dictionary("dictionaries/paper-color")

	if err != nil {
		return factory.Options{}, err
	}

	paperFactureDictionary, err := opts.Translator.Dictionary("dictionaries/paper-facture")

	if err != nil {
		return factory.Options{}, err
	}

	printFormatDictionary, err := opts.Translator.Dictionary("dictionaries/print-format")

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

		UnitLaminateType: factory.UnitLaminateTypeOptions{
			Dictionary: laminateTypeDictionary,
		},

		UnitPaperColor: factory.UnitPaperColorOptions{
			Dictionary: paperColorDictionary,
		},

		UnitPaperFacture: factory.UnitPaperFactureOptions{
			Dictionary: paperFactureDictionary,
		},

		UnitPrintFormatFacture: factory.UnitPrintFormatOptions{
			Dictionary: printFormatDictionary,
		},
	}, nil
}

func NewDictionariesLaminateTypeAPI(ctx context.Context, opts modules.Options) (*usecase_api.LaminateType, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init dictionaries laminate type API")

	return factory_api.NewLaminateType(opts.PostgresAdapter, opts.UsecaseHelper), nil
}

func NewDictionariesPaperColorAPI(ctx context.Context, opts modules.Options) (*usecase_api.PaperColor, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init dictionaries paper color API")

	return factory_api.NewPaperColor(opts.PostgresAdapter, opts.UsecaseHelper), nil
}

func NewDictionariesPaperFactureAPI(ctx context.Context, opts modules.Options) (*usecase_api.PaperFacture, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init dictionaries paper facture API")

	return factory_api.NewPaperFacture(opts.PostgresAdapter, opts.UsecaseHelper), nil
}

func NewDictionariesPrintFormatAPI(ctx context.Context, opts modules.Options) (*usecase_api.PrintFormat, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init dictionaries print format API")

	return factory_api.NewPrintFormat(opts.PostgresAdapter, opts.UsecaseHelper), nil
}
