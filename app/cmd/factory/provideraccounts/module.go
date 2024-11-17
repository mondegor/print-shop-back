package provideraccounts

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore/mrapp"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/factory/provideraccounts"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/shared/validate"
	pkgvalidate "github.com/mondegor/print-shop-back/pkg/provideraccounts/validate"
)

// NewModuleOptions - создаёт объект provideraccounts.Options.
func NewModuleOptions(_ context.Context, opts app.Options) (provideraccounts.Options, error) {
	fileAPI, err := opts.FileProviderPool.ProviderAPI(
		opts.Cfg.ModulesSettings.ProviderAccount.CompanyPageLogo.FileProvider,
	)
	if err != nil {
		return provideraccounts.Options{}, err
	}

	return provideraccounts.Options{
		EventEmitter:  opts.EventEmitter,
		UseCaseHelper: mrapp.NewUseCaseErrorWrapper(),
		DBConnManager: opts.PostgresConnManager,
		Locker:        opts.Locker,
		RequestParsers: provideraccounts.RequestParsers{
			// Parser:       opts.RequestParsers.Parser,
			// ExtendParser: opts.RequestParsers.ExtendParser,
			ModuleParser: validate.NewParser(
				opts.RequestParsers.ExtendParser,
				opts.RequestParsers.ImageLogo,
				pkgvalidate.NewPublicStatusParser(),
			),
		},
		ResponseSender: opts.ResponseSenders.Sender,

		UnitCompanyPage: provideraccounts.UnitCompanyPageOptions{
			LogoFileAPI:    fileAPI,
			LogoURLBuilder: opts.ImageURLBuilder,
		},

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}
