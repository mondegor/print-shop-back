package factory

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/app"
	calculationsalgopub "github.com/mondegor/print-shop-back/internal/factory/calculations/algo/section/pub"
	calculationsquerypub "github.com/mondegor/print-shop-back/internal/factory/calculations/queryhistory/section/pub"
	catalogboxadm "github.com/mondegor/print-shop-back/internal/factory/catalog/box/section/adm"
	catalogboxpub "github.com/mondegor/print-shop-back/internal/factory/catalog/box/section/pub"
	cataloglaminateadm "github.com/mondegor/print-shop-back/internal/factory/catalog/laminate/section/adm"
	cataloglaminatepub "github.com/mondegor/print-shop-back/internal/factory/catalog/laminate/section/pub"
	catalogpaperadm "github.com/mondegor/print-shop-back/internal/factory/catalog/paper/section/adm"
	catalogpaperpub "github.com/mondegor/print-shop-back/internal/factory/catalog/paper/section/pub"
	controlselementtemplateadm "github.com/mondegor/print-shop-back/internal/factory/controls/elementtemplate/section/adm"
	controlssubmitformadm "github.com/mondegor/print-shop-back/internal/factory/controls/submitform/section/adm"
	controlssubmitformpub "github.com/mondegor/print-shop-back/internal/factory/controls/submitform/section/pub"
	dictionariesmaterialtypeadm "github.com/mondegor/print-shop-back/internal/factory/dictionaries/materialtype/section/adm"
	dictionariesmaterialtypepub "github.com/mondegor/print-shop-back/internal/factory/dictionaries/materialtype/section/pub"
	dictionariespapercoloradm "github.com/mondegor/print-shop-back/internal/factory/dictionaries/papercolor/section/adm"
	dictionariespapercolorpub "github.com/mondegor/print-shop-back/internal/factory/dictionaries/papercolor/section/pub"
	dictionariespaperfactureadm "github.com/mondegor/print-shop-back/internal/factory/dictionaries/paperfacture/section/adm"
	dictionariespaperfacturepub "github.com/mondegor/print-shop-back/internal/factory/dictionaries/paperfacture/section/pub"
	dictionariesprintformatadm "github.com/mondegor/print-shop-back/internal/factory/dictionaries/printformat/section/adm"
	dictionariesprintformatpub "github.com/mondegor/print-shop-back/internal/factory/dictionaries/printformat/section/pub"
	filestationpub "github.com/mondegor/print-shop-back/internal/factory/filestation/section/pub"
	provideraccountsadm "github.com/mondegor/print-shop-back/internal/factory/provideraccounts/section/adm"
	provideraccountsprov "github.com/mondegor/print-shop-back/internal/factory/provideraccounts/section/prov"
	provideraccountpub "github.com/mondegor/print-shop-back/internal/factory/provideraccounts/section/pub"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
)

const (
	restServerCaption = "RestServer"
)

// NewRestServer - создаёт объект mrserver.ServerAdapter.
func NewRestServer(ctx context.Context, opts app.Options) (*mrserver.ServerAdapter, error) {
	mrlog.Ctx(ctx).Info().Msgf("Create and init '%s'", restServerCaption)

	router, err := NewRestRouter(ctx, opts, opts.Translator)
	if err != nil {
		return nil, err
	}

	// section: admin-api
	{
		sectionAdminAPI := NewAppSectionAdminAPI(ctx, opts)

		if err = RegisterSystemHandlers(ctx, opts.Cfg, router, sectionAdminAPI); err != nil {
			return nil, err
		}

		registerControllersFunc := registerControllers(
			router,
			mrfactory.WithMiddlewareCheckAccess(ctx, sectionAdminAPI, opts.AccessControl),
		)

		for _, createFunc := range getAdminAPIControllers(ctx, opts) {
			list, err := createFunc()
			if err != nil {
				return nil, err
			}

			registerControllersFunc(list)
		}
	}

	// section: providers-api
	{
		sectionProvidersAPI := NewAppSectionProvidersAPI(ctx, opts)

		if err = RegisterSystemHandlers(ctx, opts.Cfg, router, sectionProvidersAPI); err != nil {
			return nil, err
		}

		registerControllersFunc := registerControllers(
			router,
			mrfactory.WithMiddlewareCheckAccess(ctx, sectionProvidersAPI, opts.AccessControl),
		)

		for _, createFunc := range getProvidersAPIControllers(ctx, opts) {
			list, err := createFunc()
			if err != nil {
				return nil, err
			}

			registerControllersFunc(list)
		}
	}

	// section: public
	{
		sectionPublicAPI := NewAppSectionPublicAPI(ctx, opts)

		if err = RegisterSystemHandlers(ctx, opts.Cfg, router, sectionPublicAPI); err != nil {
			return nil, err
		}

		registerControllersFunc := registerControllers(
			router,
			mrfactory.WithMiddlewareCheckAccess(ctx, sectionPublicAPI, opts.AccessControl),
		)

		for _, createFunc := range getPublicAPIControllers(ctx, opts) {
			list, err := createFunc()
			if err != nil {
				return nil, err
			}

			registerControllersFunc(list)
		}
	}

	srvOpts := opts.Cfg.Servers.RestServer

	return mrserver.NewServerAdapter(
		ctx,
		mrserver.ServerOptions{
			Caption:         restServerCaption,
			Handler:         router,
			ReadTimeout:     srvOpts.ReadTimeout,
			WriteTimeout:    srvOpts.WriteTimeout,
			ShutdownTimeout: srvOpts.ShutdownTimeout,
			Listen: mrserver.ListenOptions{
				BindIP: srvOpts.Listen.BindIP,
				Port:   srvOpts.Listen.Port,
			},
		},
	), nil
}

func registerControllers(router mrserver.HttpRouter, operations ...mrfactory.PrepareHandlerFunc) func(list []mrserver.HttpController) {
	return func(list []mrserver.HttpController) {
		router.Register(
			mrfactory.PrepareEachController(list, operations...)...,
		)
	}
}

func getAdminAPIControllers(ctx context.Context, opts app.Options) []func() (list []mrserver.HttpController, err error) {
	return []func() (list []mrserver.HttpController, err error){
		func() ([]mrserver.HttpController, error) {
			return catalogboxadm.CreateModule(ctx, opts.CatalogBoxModule)
		},
		func() ([]mrserver.HttpController, error) {
			return cataloglaminateadm.CreateModule(ctx, opts.CatalogLaminateModule)
		},
		func() ([]mrserver.HttpController, error) {
			return catalogpaperadm.CreateModule(ctx, opts.CatalogPaperModule)
		},
		func() ([]mrserver.HttpController, error) {
			return controlselementtemplateadm.CreateModule(ctx, opts.ControlsElementTemplateModule)
		},
		func() ([]mrserver.HttpController, error) {
			return controlssubmitformadm.CreateModule(ctx, opts.ControlsSubmitFormModule)
		},
		func() ([]mrserver.HttpController, error) {
			return dictionariesmaterialtypeadm.CreateModule(ctx, opts.DictionariesMaterialTypeModule)
		},
		func() ([]mrserver.HttpController, error) {
			return dictionariespapercoloradm.CreateModule(ctx, opts.DictionariesPaperColorModule)
		},
		func() ([]mrserver.HttpController, error) {
			return dictionariespaperfactureadm.CreateModule(ctx, opts.DictionariesPaperFactureModule)
		},
		func() ([]mrserver.HttpController, error) {
			return dictionariesprintformatadm.CreateModule(ctx, opts.DictionariesPrintFormatModule)
		},
		func() ([]mrserver.HttpController, error) {
			return provideraccountsadm.CreateModule(ctx, opts.ProviderAccountsModule)
		},
	}
}

func getProvidersAPIControllers(ctx context.Context, opts app.Options) []func() (list []mrserver.HttpController, err error) {
	return []func() (list []mrserver.HttpController, err error){
		func() ([]mrserver.HttpController, error) {
			return provideraccountsprov.CreateModule(ctx, opts.ProviderAccountsModule)
		},
	}
}

func getPublicAPIControllers(ctx context.Context, opts app.Options) []func() (list []mrserver.HttpController, err error) {
	return []func() (list []mrserver.HttpController, err error){
		func() ([]mrserver.HttpController, error) {
			return calculationsquerypub.CreateModule(ctx, opts.CalculationsQueryHistoryModule)
		},
		func() ([]mrserver.HttpController, error) {
			return calculationsalgopub.CreateModule(ctx, opts.CalculationsAlgoModule)
		},
		func() ([]mrserver.HttpController, error) {
			return catalogboxpub.CreateModule(ctx, opts.CatalogBoxModule)
		},
		func() ([]mrserver.HttpController, error) {
			return cataloglaminatepub.CreateModule(ctx, opts.CatalogLaminateModule)
		},
		func() ([]mrserver.HttpController, error) {
			return catalogpaperpub.CreateModule(ctx, opts.CatalogPaperModule)
		},
		func() ([]mrserver.HttpController, error) {
			return controlssubmitformpub.CreateModule(ctx, opts.ControlsSubmitFormModule)
		},
		func() ([]mrserver.HttpController, error) {
			return dictionariesmaterialtypepub.CreateModule(ctx, opts.DictionariesMaterialTypeModule)
		},
		func() ([]mrserver.HttpController, error) {
			return dictionariespapercolorpub.CreateModule(ctx, opts.DictionariesPaperColorModule)
		},
		func() ([]mrserver.HttpController, error) {
			return dictionariespaperfacturepub.CreateModule(ctx, opts.DictionariesPaperFactureModule)
		},
		func() ([]mrserver.HttpController, error) {
			return dictionariesprintformatpub.CreateModule(ctx, opts.DictionariesPrintFormatModule)
		},
		func() ([]mrserver.HttpController, error) {
			return filestationpub.CreateModule(ctx, opts.FileStationModule)
		},
		func() ([]mrserver.HttpController, error) {
			return provideraccountpub.CreateModule(ctx, opts.ProviderAccountsModule)
		},
	}
}
