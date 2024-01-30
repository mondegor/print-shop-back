package factory

import (
	"context"
	"print-shop-back/internal/modules"
	view_shared "print-shop-back/internal/modules/provider-accounts/controller/http_v1/shared/view"
	"print-shop-back/internal/modules/provider-accounts/factory"
)

func NewProviderAccountsModuleOptions(ctx context.Context, opts modules.Options) (factory.Options, error) {
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
				opts.RequestParsers.SortPage,
				opts.RequestParsers.String,
				opts.RequestParsers.Validator,
			),
		},
		ResponseSender: opts.ResponseSender,

		UnitCompanyPage: factory.UnitCompanyPageOptions{
			LogoFileAPI:    fileAPI,
			LogoURLBuilder: NewBuilderImagesURL(opts.Cfg),
		},
	}, nil
}
