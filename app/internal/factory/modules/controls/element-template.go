package factory_controls

import (
	"context"
	"print-shop-back/internal"
	view_shared "print-shop-back/internal/modules/controls/element-template/controller/http_v1/shared/view"
	"print-shop-back/internal/modules/controls/element-template/factory"
	"print-shop-back/pkg/modules/controls/view"
)

func NewElementTemplateModuleOptions(ctx context.Context, opts app.Options) (factory.Options, error) {
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
			opts.RequestParsers.FileJson,
			view.NewDetailingParser(),
		),
		ResponseSender: opts.ResponseSender,

		UnitElementTemplate: factory.UnitElementTemplateOptions{},

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}
