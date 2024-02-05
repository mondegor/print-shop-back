package factory_dictionaries

import (
	"context"
	"print-shop-back/internal"
	view_shared "print-shop-back/internal/modules/dictionaries/laminate-type/controller/http_v1/shared/view"
	"print-shop-back/internal/modules/dictionaries/laminate-type/factory"
	factory_api "print-shop-back/internal/modules/dictionaries/laminate-type/factory/api"
	usecase_api "print-shop-back/internal/modules/dictionaries/laminate-type/usecase/api"

	"github.com/mondegor/go-webcore/mrlog"
)

func NewLaminateTypeModuleOptions(ctx context.Context, opts app.Options) (factory.Options, error) {
	laminateTypeDictionary, err := opts.Translator.Dictionary("dictionaries/laminate-types")

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
	}, nil
}

func NewLaminateTypeAPI(ctx context.Context, opts app.Options) (*usecase_api.LaminateType, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init dictionaries laminate type API")

	return factory_api.NewLaminateType(opts.PostgresAdapter, opts.UsecaseHelper), nil
}
