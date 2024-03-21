package factory_provideraccounts

import (
	"context"
	"print-shop-back/internal"
	view_shared "print-shop-back/internal/modules/provider-accounts/controller/http_v1/shared/view"
	"print-shop-back/internal/modules/provider-accounts/factory"
	view_shared2 "print-shop-back/pkg/modules/provider-accounts/view"
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
		RequestParser: view_shared.NewParser(
			opts.RequestParsers.Int64,
			opts.RequestParsers.ItemStatus,
			opts.RequestParsers.KeyInt32,
			opts.RequestParsers.ListSorter,
			opts.RequestParsers.ListPager,
			opts.RequestParsers.String,
			opts.RequestParsers.Validator,
			opts.RequestParsers.ImageLogo,
			view_shared2.NewPublicStatusParser(),
		),
		ResponseSender: opts.ResponseSender,

		UnitCompanyPage: factory.UnitCompanyPageOptions{
			LogoFileAPI:    fileAPI,
			LogoURLBuilder: opts.ImageURLBuilder,
		},

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}
