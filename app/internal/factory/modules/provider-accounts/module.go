package factory_provideraccounts

import (
	"context"
	"print-shop-back/internal"
	view_shared "print-shop-back/internal/modules/provider-accounts/controller/http_v1/shared/view"
	"print-shop-back/internal/modules/provider-accounts/factory"
)

func NewModuleOptions(ctx context.Context, opts app.Options) (factory.Options, error) {
	fileAPI, err := opts.FileProviderPool.Provider(
		opts.Cfg.ModulesSettings.ProviderAccount.CompanyPageLogo.FileProvider,
	)

	if err != nil {
		return factory.Options{}, err
	}

	return factory.Options{
		EventEmitter:    opts.EventEmitter,
		UsecaseHelper:   opts.UsecaseHelper,
		PostgresAdapter: opts.PostgresAdapter,
		Locker:          opts.Locker,
		RequestParsers: factory.RequestParsers{
			String: opts.RequestParsers.String,
			Image:  opts.RequestParsers.Image,
			Parser: view_shared.NewParser(
				opts.RequestParsers.Int64,
				opts.RequestParsers.ItemStatus,
				opts.RequestParsers.KeyInt32,
				opts.RequestParsers.ListSorter,
				opts.RequestParsers.ListPager,
				opts.RequestParsers.String,
				opts.RequestParsers.Validator,
			),
		},
		ResponseSender: opts.ResponseSender,

		UnitCompanyPage: factory.UnitCompanyPageOptions{
			LogoFileAPI:    fileAPI,
			LogoURLBuilder: opts.ImageURLBuilder,
		},

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}
