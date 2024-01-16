package factory

import (
	"print-shop-back/internal/modules"
	"print-shop-back/internal/modules/controls/factory"
	factory_api "print-shop-back/internal/modules/controls/factory/api"
)

func NewControlsOptions(opts *modules.Options) (*factory.Options, error) {
	return &factory.Options{
		Logger:          opts.Logger,
		EventBox:        opts.EventBox,
		ServiceHelper:   opts.ServiceHelper,
		PostgresAdapter: opts.PostgresAdapter,
		OrdererAPI:      opts.OrdererAPI,

		ElementTemplateAPI: factory_api.NewElementTemplate(opts.PostgresAdapter, opts.ServiceHelper),

		UnitElementTemplate: &factory.UnitElementTemplateOptions{},
	}, nil
}
