package rest

import (
	"net/http"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mraccess"
	"github.com/mondegor/go-webcore/mraccess/section"
	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresp"

	"github.com/mondegor/print-shop-back/internal/app"
	provideraccounts "github.com/mondegor/print-shop-back/internal/factory/provideraccounts/section/prov"
	"github.com/mondegor/print-shop-back/internal/initing"
	provideraccountsvalidate "github.com/mondegor/print-shop-back/internal/provideraccounts/shared/validate"
	pkgprovideraccountsvalidate "github.com/mondegor/print-shop-back/pkg/provideraccounts/validate"
)

// RegisterRestRouterProvHandlers - регистрирует в указанном роутере обработчики секции ProvidersAPI.
func RegisterRestRouterProvHandlers(router mrserver.HttpRouter, opts app.Options, sect *section.RoutingSection, memberProvider mraccess.MemberProvider) error {
	router.HandlerFunc(http.MethodGet, sect.BuildPath("/"), mrresp.HandlerGetStatusOkAsJSON(opts.Logger))
	prepareHandler := mrinit.WithMiddlewareCheckAccess(opts.Logger, sect, memberProvider, opts.RealmKindRights, opts.PermsProvider)

	controllers, err := initing.CreateHttpControllers(opts.Logger, getProvidersAPIControllers(opts), prepareHandler)
	if err != nil {
		return err
	}

	router.Register(controllers...)

	return nil
}

func getProvidersAPIControllers(opts app.Options) []initing.HttpModule {
	return []initing.HttpModule{
		provideraccounts.InitHttpModule(
			opts.Logger,
			opts.EventEmitter,
			opts.UseCaseErrorWrapper,
			opts.ImageUserErrorWrapper,
			opts.PostgresConnManager,
			opts.Locker,
			provideraccountsvalidate.NewParser(
				opts.RequestParsers.ExtendParser,
				opts.RequestParsers.User,
				opts.RequestParsers.ImageLogo,
				pkgprovideraccountsvalidate.NewPublicStatusParser(opts.Logger),
			),
			opts.ResponseSenders.Sender,
			func() (mrstorage.FileProviderAPI, error) {
				return opts.FileProviderPool.ProviderAPI(
					opts.Cfg.ModulesSettings.ProviderAccount.CompanyPageLogo.FileProvider,
				)
			},
			opts.ImageURLBuilder,
		),
	}
}
