package factory_controls

import (
	"context"
	"print-shop-back/internal"
	factory_api "print-shop-back/internal/modules/controls/element-template/factory/api"
	view_shared "print-shop-back/internal/modules/controls/submit-form/controller/http_v1/shared/view"
	"print-shop-back/internal/modules/controls/submit-form/factory"
	"print-shop-back/pkg/modules/controls/view"
)

func NewSubmitFormModuleOptions(ctx context.Context, opts app.Options) (factory.Options, error) {
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
			opts.RequestParsers.UUID,
			opts.RequestParsers.Validator,
			opts.RequestParsers.FileJson,
			view.NewDetailingParser(),
		),
		ResponseSender: opts.ResponseSender,

		ElementTemplateAPI: factory_api.NewElementTemplate(opts.PostgresAdapter, opts.UsecaseHelper),
		OrdererAPI:         opts.OrdererAPI,

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}
