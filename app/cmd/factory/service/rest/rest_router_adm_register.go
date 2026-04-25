package rest

import (
	"net/http"

	"github.com/mondegor/go-webcore/mraccess"
	"github.com/mondegor/go-webcore/mrcore/initing"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresp"

	"github.com/mondegor/print-shop-back/internal/app"
	controlselementtemplatevalidate "github.com/mondegor/print-shop-back/internal/controls/elementtemplate/shared/validate"
	controlssubmitformvalidate "github.com/mondegor/print-shop-back/internal/controls/submitform/shared/validate"
	catalogbox "github.com/mondegor/print-shop-back/internal/factory/catalog/box/section/adm"
	cataloglaminate "github.com/mondegor/print-shop-back/internal/factory/catalog/laminate/section/adm"
	catalogpaper "github.com/mondegor/print-shop-back/internal/factory/catalog/paper/section/adm"
	controlssubmitformapi "github.com/mondegor/print-shop-back/internal/factory/controls/elementtemplate/api/header"
	controlselementtemplate "github.com/mondegor/print-shop-back/internal/factory/controls/elementtemplate/section/adm"
	controlssubmitform "github.com/mondegor/print-shop-back/internal/factory/controls/submitform/section/adm"
	dictionariesmaterialtype "github.com/mondegor/print-shop-back/internal/factory/dictionaries/materialtype/section/adm"
	dictionariespapercolor "github.com/mondegor/print-shop-back/internal/factory/dictionaries/papercolor/section/adm"
	dictionariespaperfacture "github.com/mondegor/print-shop-back/internal/factory/dictionaries/paperfacture/section/adm"
	dictionariesprintformat "github.com/mondegor/print-shop-back/internal/factory/dictionaries/printformat/section/adm"
	provideraccounts "github.com/mondegor/print-shop-back/internal/factory/provideraccounts/section/adm"
	provideraccountsvalidate "github.com/mondegor/print-shop-back/internal/provideraccounts/shared/validate"
	pkgcontrolsvalidate "github.com/mondegor/print-shop-back/pkg/controls/validate"
	pkgprovideraccountsvalidate "github.com/mondegor/print-shop-back/pkg/provideraccounts/validate"
)

// RegisterRestRouterAdmHandlers - регистрирует в указанном роутере обработчики секции AdminAPI.
func RegisterRestRouterAdmHandlers(
	router mrserver.HttpRouter,
	opts app.Options,
	actionGroup *mraccess.ActionGroup,
	userProvider mraccess.UserProvider,
) error {
	router.HandlerFunc(http.MethodGet, actionGroup.BasePath.BuildPath("/"), mrresp.HandlerGetStatusOkAsJSON(opts.Logger))

	controllers, err := initing.CreateHttpControllers(
		opts.Logger,
		getAdminAPIControllers(opts),
		initing.WithCheckAccessMiddleware(opts.Logger, actionGroup, userProvider, opts.PermsProvider),
	)
	if err != nil {
		return err
	}

	router.Register(controllers...)

	return nil
}

func getAdminAPIControllers(opts app.Options) []initing.HttpModule {
	return []initing.HttpModule{
		catalogbox.InitHttpModule(
			opts.Logger,
			opts.EventEmitter,
			opts.PostgresConnManager,
			opts.RequestParsers.ExtendParser,
			opts.ResponseSenders.Sender,
			int(opts.Cfg.ModuleSettings.General.PageSizeMax),
		),
		cataloglaminate.InitHttpModule(
			opts.Logger,
			opts.EventEmitter,
			opts.PostgresConnManager,
			opts.RequestParsers.ExtendParser,
			opts.ResponseSenders.Sender,
			opts.DictionariesMaterialTypeAPI,
			int(opts.Cfg.ModuleSettings.General.PageSizeMax),
		),
		catalogpaper.InitHttpModule(
			opts.Logger,
			opts.EventEmitter,
			opts.PostgresConnManager,
			opts.RequestParsers.ExtendParser,
			opts.ResponseSenders.Sender,
			opts.DictionariesMaterialTypeAPI,
			opts.DictionariesPaperColorAPI,
			opts.DictionariesPaperFactureAPI,
			int(opts.Cfg.ModuleSettings.General.PageSizeMax),
		),
		controlselementtemplate.InitHttpModule(
			opts.Logger,
			opts.EventEmitter,
			opts.PostgresConnManager,
			controlselementtemplatevalidate.NewParser(
				opts.RequestParsers.ExtendParser,
				opts.RequestParsers.FileJson,
				pkgcontrolsvalidate.NewDetailingParser(opts.Logger),
			),
			opts.ResponseSenders.FileSender,
			int(opts.Cfg.ModuleSettings.General.PageSizeMax),
		),
		controlssubmitform.InitHttpModule(
			opts.Logger,
			opts.EventEmitter,
			opts.PostgresConnManager,
			opts.Locker,
			controlssubmitformvalidate.NewParser(
				opts.RequestParsers.ExtendParser,
				opts.RequestParsers.FileJson,
				pkgcontrolsvalidate.NewDetailingParser(opts.Logger),
			),
			opts.ResponseSenders.Sender,
			opts.ResponseSenders.FileSender,
			controlssubmitformapi.NewElementTemplate(
				opts.PostgresConnManager,
				opts.Tracer,
			),
			int(opts.Cfg.ModuleSettings.General.PageSizeMax),
		),
		dictionariesmaterialtype.InitHttpModule(
			opts.Logger,
			opts.EventEmitter,
			opts.PostgresConnManager,
			opts.RequestParsers.ExtendParser,
			opts.ResponseSenders.Sender,
			int(opts.Cfg.ModuleSettings.General.PageSizeMax),
		),
		dictionariespapercolor.InitHttpModule(
			opts.Logger,
			opts.EventEmitter,
			opts.PostgresConnManager,
			opts.RequestParsers.ExtendParser,
			opts.ResponseSenders.Sender,
			int(opts.Cfg.ModuleSettings.General.PageSizeMax),
		),
		dictionariespaperfacture.InitHttpModule(
			opts.Logger,
			opts.EventEmitter,
			opts.PostgresConnManager,
			opts.RequestParsers.ExtendParser,
			opts.ResponseSenders.Sender,
			int(opts.Cfg.ModuleSettings.General.PageSizeMax),
		),
		dictionariesprintformat.InitHttpModule(
			opts.Logger,
			opts.EventEmitter,
			opts.PostgresConnManager,
			opts.RequestParsers.ExtendParser,
			opts.ResponseSenders.Sender,
			int(opts.Cfg.ModuleSettings.General.PageSizeMax),
		),
		provideraccounts.InitHttpModule(
			opts.Logger,
			opts.PostgresConnManager,
			provideraccountsvalidate.NewParser(
				opts.RequestParsers.ExtendParser,
				opts.RequestParsers.User,
				opts.RequestParsers.ImageLogo,
				pkgprovideraccountsvalidate.NewPublicStatusParser(opts.Logger),
			),
			opts.ResponseSenders.Sender,
			opts.ImageURLBuilder,
			int(opts.Cfg.ModuleSettings.General.PageSizeMax),
		),
	}
}
