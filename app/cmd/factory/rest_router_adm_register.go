package factory

import (
	"context"
	"net/http"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresp"

	"github.com/mondegor/print-shop-back/internal/app"
	catalogboxadm "github.com/mondegor/print-shop-back/internal/factory/catalog/box/section/adm"
	cataloglaminateadm "github.com/mondegor/print-shop-back/internal/factory/catalog/laminate/section/adm"
	catalogpaperadm "github.com/mondegor/print-shop-back/internal/factory/catalog/paper/section/adm"
	controlselementtemplateadm "github.com/mondegor/print-shop-back/internal/factory/controls/elementtemplate/section/adm"
	controlssubmitformadm "github.com/mondegor/print-shop-back/internal/factory/controls/submitform/section/adm"
	dictionariesmaterialtypeadm "github.com/mondegor/print-shop-back/internal/factory/dictionaries/materialtype/section/adm"
	dictionariespapercoloradm "github.com/mondegor/print-shop-back/internal/factory/dictionaries/papercolor/section/adm"
	dictionariespaperfactureadm "github.com/mondegor/print-shop-back/internal/factory/dictionaries/paperfacture/section/adm"
	dictionariesprintformatadm "github.com/mondegor/print-shop-back/internal/factory/dictionaries/printformat/section/adm"
	provideraccountsadm "github.com/mondegor/print-shop-back/internal/factory/provideraccounts/section/adm"
)

// RegisterRestRouterAdmHandlers - регистрирует в указанном роутере обработчики секции AdminAPI.
func RegisterRestRouterAdmHandlers(ctx context.Context, router mrserver.HttpRouter, opts app.Options) error {
	section := NewAppSectionAdminAPI(ctx, opts)
	prepareHandler := mrfactory.WithMiddlewareCheckAccess(ctx, section, opts.AccessControl)
	router.HandlerFunc(http.MethodGet, section.BuildPath("/"), mrresp.HandlerGetStatusOkAsJSON())

	for _, createFunc := range getAdminAPIControllers(ctx, opts) {
		list, err := createFunc()
		if err != nil {
			return err
		}

		router.Register(
			mrfactory.PrepareEachController(list, prepareHandler)...,
		)
	}

	return nil
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
