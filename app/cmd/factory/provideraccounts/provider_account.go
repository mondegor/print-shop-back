package provideraccounts

import (
	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/factory/provideraccounts"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/shared/validate"
	pkgvalidate "github.com/mondegor/print-shop-back/pkg/provideraccounts/validate"
)

// NewModuleOptions - создаёт объект provideraccounts.Options.
func NewModuleOptions(opts app.Options) (provideraccounts.Options, error) {
	fileAPI, err := opts.FileProviderPool.ProviderAPI(
		opts.Cfg.ModulesSettings.ProviderAccount.CompanyPageLogo.FileProvider,
	)
	if err != nil {
		return provideraccounts.Options{}, err
	}

	return provideraccounts.Options{
		Logger:                opts.Logger,
		EventEmitter:          opts.EventEmitter,
		UsecaseErrorWrapper:   opts.UsecaseErrorWrapper,
		ImageUserErrorWrapper: opts.ImageUserErrorWrapper,
		DBConnManager:         opts.PostgresConnManager,
		Locker:                opts.Locker,
		RequestParsers: provideraccounts.RequestParsers{
			// Parser:       opts.RequestParsers.Parser,
			// ExtendParser: opts.RequestParsers.ExtendParser,
			ModuleParser: validate.NewParser(
				opts.RequestParsers.ExtendParser,
				opts.RequestParsers.User,
				opts.RequestParsers.ImageLogo,
				pkgvalidate.NewPublicStatusParser(opts.Logger),
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
