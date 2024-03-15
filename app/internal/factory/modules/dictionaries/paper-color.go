package factory_dictionaries

import (
	"context"
	"print-shop-back/internal"
	view_shared "print-shop-back/internal/modules/dictionaries/paper-color/controller/http_v1/shared/view"
	"print-shop-back/internal/modules/dictionaries/paper-color/factory"
	factory_api "print-shop-back/internal/modules/dictionaries/paper-color/factory/api"
	usecase_api "print-shop-back/internal/modules/dictionaries/paper-color/usecase/api"

	"github.com/mondegor/go-webcore/mrlog"
)

func NewPaperColorModuleOptions(ctx context.Context, opts app.Options) (factory.Options, error) {
	paperColorDictionary, err := opts.Translator.Dictionary("dictionaries/paper-colors")

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

		UnitPaperColor: factory.UnitPaperColorOptions{
			Dictionary: paperColorDictionary,
		},

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}

func NewPaperColorAPI(ctx context.Context, opts app.Options) (*usecase_api.PaperColor, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init dictionaries paper color API")

	return factory_api.NewPaperColor(opts.PostgresAdapter, opts.UsecaseHelper), nil
}
