package rest

import (
	"net/http"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mraccess"
	"github.com/mondegor/go-webcore/mraccess/section"
	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrpath"
	"github.com/mondegor/go-webcore/mrpath/placeholderpath"
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
	"github.com/mondegor/print-shop-back/internal/initing"
	provideraccountsvalidate "github.com/mondegor/print-shop-back/internal/provideraccounts/shared/validate"
	pkgcontrolsvalidate "github.com/mondegor/print-shop-back/pkg/controls/validate"
	pkgprovideraccountsvalidate "github.com/mondegor/print-shop-back/pkg/provideraccounts/validate"
)

// RegisterRestRouterPubHandlers - регистрирует в указанном роутере обработчики секции PublicAPI.
func RegisterRestRouterPubHandlers(router mrserver.HttpRouter, opts app.Options, sect *section.RoutingSection, memberProvider mraccess.MemberProvider) error {
	router.HandlerFunc(http.MethodGet, sect.BuildPath("/"), mrresp.HandlerGetStatusOkAsJSON(opts.Logger))
	prepareHandler := mrinit.WithMiddlewareCheckAccess(opts.Logger, sect, memberProvider, opts.RealmKindRights, opts.PermsProvider)

	controllers, err := initing.CreateHttpControllers(opts.Logger, getPublicAPIControllers(opts), prepareHandler)
	if err != nil {
		return err
	}

	router.Register(controllers...)

	return nil
}

func getPublicAPIControllers(opts app.Options) []initing.HttpModule {
	return []initing.HttpModule{
		catalogbox.InitHttpModule(
			opts.UseCaseErrorWrapper,
			opts.PostgresConnManager,
			opts.RequestParsers.Parser,
			opts.ResponseSenders.Sender,
		),
		cataloglaminate.InitHttpModule(
			opts.UseCaseErrorWrapper,
			opts.PostgresConnManager,
			opts.RequestParsers.Parser,
			opts.ResponseSenders.Sender,
		),
		catalogpaper.InitHttpModule(
			opts.UseCaseErrorWrapper,
			opts.PostgresConnManager,
			opts.RequestParsers.Parser,
			opts.ResponseSenders.Sender,
		),
		controlssubmitform.InitHttpModule(
			opts.UseCaseErrorWrapper,
			opts.PostgresConnManager,
			controlssubmitformvalidate.NewParser(
				opts.RequestParsers.ExtendParser,
				opts.RequestParsers.FileJson,
				pkgcontrolsvalidate.NewDetailingParser(opts.Logger),
			),
			opts.ResponseSenders.Sender,
		),
		dictionariesmaterialtype.InitHttpModule(
			opts.UseCaseErrorWrapper,
			opts.PostgresConnManager,
			opts.RequestParsers.Parser,
			opts.ResponseSenders.Sender,
		),
		dictionariespapercolor.InitHttpModule(
			opts.UseCaseErrorWrapper,
			opts.PostgresConnManager,
			opts.RequestParsers.Parser,
			opts.ResponseSenders.Sender,
		),
		dictionariespaperfacture.InitHttpModule(
			opts.UseCaseErrorWrapper,
			opts.PostgresConnManager,
			opts.RequestParsers.Parser,
			opts.ResponseSenders.Sender,
		),
		dictionariesprintformat.InitHttpModule(
			opts.UseCaseErrorWrapper,
			opts.PostgresConnManager,
			opts.RequestParsers.Parser,
			opts.ResponseSenders.Sender,
		),
		provideraccount.InitHttpModule(
			opts.UseCaseErrorWrapper,
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
			opts.UseCaseErrorWrapper,
			opts.RequestParsers.String,
			opts.ResponseSenders.FileSender,
			func() (mrstorage.FileProviderAPI, mrpath.PathBuilder, error) {
				fileAPI, err := opts.FileProviderPool.ProviderAPI(
					opts.Cfg.ModulesSettings.FileStation.ImageProxy.FileProvider,
				)
				if err != nil {
					return nil, nil, err
				}

				basePath, err := placeholderpath.New(
					opts.Cfg.ModulesSettings.FileStation.ImageProxy.BasePath,
					placeholderpath.Placeholder,
				)
				if err != nil {
					return nil, nil, err
				}

				return fileAPI, basePath, nil
			},
		),
		calculationsquery.InitHttpModule(
			opts.EventEmitter,
			opts.UseCaseErrorWrapper,
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
