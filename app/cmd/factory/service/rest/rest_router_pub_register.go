package rest

import (
	"net/http"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrpath"
	"github.com/mondegor/go-webcore/mraccess"
	"github.com/mondegor/go-webcore/mrcore/initing"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresp"

	"github.com/mondegor/print-shop-back/internal/app"
	controlssubmitformvalidate "github.com/mondegor/print-shop-back/internal/controls/submitform/shared/validate"
	calculationsalgo "github.com/mondegor/print-shop-back/internal/factory/calculations/algo/section/pub"
	calculationsquery "github.com/mondegor/print-shop-back/internal/factory/calculations/queryhistory/section/pub"
	catalogbox "github.com/mondegor/print-shop-back/internal/factory/catalog/box/section/pub"
	cataloglaminate "github.com/mondegor/print-shop-back/internal/factory/catalog/laminate/section/pub"
	catalogpaper "github.com/mondegor/print-shop-back/internal/factory/catalog/paper/section/pub"
	controlssubmitform "github.com/mondegor/print-shop-back/internal/factory/controls/submitform/section/pub"
	dictionariesmaterialtype "github.com/mondegor/print-shop-back/internal/factory/dictionaries/materialtype/section/pub"
	dictionariespapercolor "github.com/mondegor/print-shop-back/internal/factory/dictionaries/papercolor/section/pub"
	dictionariespaperfacture "github.com/mondegor/print-shop-back/internal/factory/dictionaries/paperfacture/section/pub"
	dictionariesprintformat "github.com/mondegor/print-shop-back/internal/factory/dictionaries/printformat/section/pub"
	filestation "github.com/mondegor/print-shop-back/internal/factory/filestation/section/pub"
	provideraccount "github.com/mondegor/print-shop-back/internal/factory/provideraccounts/section/pub"
	provideraccountsvalidate "github.com/mondegor/print-shop-back/internal/provideraccounts/shared/validate"
	pkgcontrolsvalidate "github.com/mondegor/print-shop-back/pkg/controls/validate"
	pkgprovideraccountsvalidate "github.com/mondegor/print-shop-back/pkg/provideraccounts/validate"
)

// RegisterRestRouterPubHandlers - регистрирует в указанном роутере обработчики секции PublicAPI.
func RegisterRestRouterPubHandlers(
	router mrserver.HttpRouter,
	opts app.Options,
	actionGroup *mraccess.ActionGroup,
	userProvider mraccess.UserProvider,
) error {
	router.HandlerFunc(http.MethodGet, actionGroup.BasePath.BuildPath("/"), mrresp.HandlerGetStatusOkAsJSON(opts.Logger))

	controllers, err := initing.CreateHttpControllers(
		opts.Logger,
		getPublicAPIControllers(opts),
		initing.WithCheckAccessMiddleware(opts.Logger, actionGroup, userProvider, opts.PermsProvider),
	)
	if err != nil {
		return err
	}

	router.Register(controllers...)

	return nil
}

func getPublicAPIControllers(opts app.Options) []initing.HttpModule {
	return []initing.HttpModule{
		catalogbox.InitHttpModule(
			opts.PostgresConnManager,
			opts.RequestParsers.Parser,
			opts.ResponseSenders.Sender,
		),
		cataloglaminate.InitHttpModule(
			opts.PostgresConnManager,
			opts.RequestParsers.Parser,
			opts.ResponseSenders.Sender,
		),
		catalogpaper.InitHttpModule(
			opts.PostgresConnManager,
			opts.RequestParsers.Parser,
			opts.ResponseSenders.Sender,
		),
		controlssubmitform.InitHttpModule(
			opts.PostgresConnManager,
			controlssubmitformvalidate.NewParser(
				opts.RequestParsers.ExtendParser,
				opts.RequestParsers.FileJson,
				pkgcontrolsvalidate.NewDetailingParser(opts.Logger),
			),
			opts.ResponseSenders.Sender,
		),
		dictionariesmaterialtype.InitHttpModule(
			opts.PostgresConnManager,
			opts.RequestParsers.Parser,
			opts.ResponseSenders.Sender,
		),
		dictionariespapercolor.InitHttpModule(
			opts.PostgresConnManager,
			opts.RequestParsers.Parser,
			opts.ResponseSenders.Sender,
		),
		dictionariespaperfacture.InitHttpModule(
			opts.PostgresConnManager,
			opts.RequestParsers.Parser,
			opts.ResponseSenders.Sender,
		),
		dictionariesprintformat.InitHttpModule(
			opts.PostgresConnManager,
			opts.RequestParsers.Parser,
			opts.ResponseSenders.Sender,
		),
		provideraccount.InitHttpModule(
			opts.PostgresConnManager,
			provideraccountsvalidate.NewParser(
				opts.RequestParsers.ExtendParser,
				opts.RequestParsers.User,
				opts.RequestParsers.ImageLogo,
				pkgprovideraccountsvalidate.NewPublicStatusParser(opts.Logger),
			),
			opts.ResponseSenders.Sender,
			opts.ImageURLBuilder,
		),
		filestation.InitHttpModule(
			opts.RequestParsers.String,
			opts.ResponseSenders.FileSender,
			func() (mrstorage.FileProviderAPI, mrpath.Builder, error) {
				fileAPI, err := opts.FileProviderPool.ProviderAPI(
					opts.Cfg.ModulesSettings.FileStation.ImageProxy.FileProvider,
				)
				if err != nil {
					return nil, nil, err
				}

				basePath, err := mrpath.NewPlaceholder(
					opts.Cfg.ModulesSettings.FileStation.ImageProxy.BasePath,
					mrpath.Placeholder,
				)
				if err != nil {
					return nil, nil, err
				}

				return fileAPI, basePath, nil
			},
		),
		calculationsquery.InitHttpModule(
			opts.EventEmitter,
			opts.PostgresConnManager,
			opts.RequestParsers.ExtendParser,
			opts.ResponseSenders.Sender,
		),
		calculationsalgo.InitHttpModule(
			opts.Logger,
			opts.EventEmitter,
			opts.RequestParsers.Parser,
			opts.ResponseSenders.Sender,
		),
	}
}
