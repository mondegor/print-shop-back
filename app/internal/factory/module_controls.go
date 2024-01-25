package factory

import (
	"print-shop-back/internal/modules"
	view_shared "print-shop-back/internal/modules/controls/controller/http_v1/shared/view"
	"print-shop-back/internal/modules/controls/factory"
	factory_api "print-shop-back/internal/modules/controls/factory/api"
)

func NewControlsOptions(opts *modules.Options) (*factory.Options, error) {
	return &factory.Options{
		Logger:          opts.Logger,
		EventBox:        opts.EventBox,
		ServiceHelper:   opts.ServiceHelper,
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

		ElementTemplateAPI: factory_api.NewElementTemplate(opts.PostgresAdapter, opts.ServiceHelper),
		OrdererAPI:         opts.OrdererAPI,

		UnitElementTemplate: &factory.UnitElementTemplateOptions{},
	}, nil
}
