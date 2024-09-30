package factory

import (
	"context"
	"net/http"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresp"

	"github.com/mondegor/print-shop-back/internal/app"
	calculationsalgopub "github.com/mondegor/print-shop-back/internal/factory/calculations/algo/section/pub"
	calculationsquerypub "github.com/mondegor/print-shop-back/internal/factory/calculations/queryhistory/section/pub"
	catalogboxpub "github.com/mondegor/print-shop-back/internal/factory/catalog/box/section/pub"
	cataloglaminatepub "github.com/mondegor/print-shop-back/internal/factory/catalog/laminate/section/pub"
	catalogpaperpub "github.com/mondegor/print-shop-back/internal/factory/catalog/paper/section/pub"
	controlssubmitformpub "github.com/mondegor/print-shop-back/internal/factory/controls/submitform/section/pub"
	dictionariesmaterialtypepub "github.com/mondegor/print-shop-back/internal/factory/dictionaries/materialtype/section/pub"
	dictionariespapercolorpub "github.com/mondegor/print-shop-back/internal/factory/dictionaries/papercolor/section/pub"
	dictionariespaperfacturepub "github.com/mondegor/print-shop-back/internal/factory/dictionaries/paperfacture/section/pub"
	dictionariesprintformatpub "github.com/mondegor/print-shop-back/internal/factory/dictionaries/printformat/section/pub"
	filestationpub "github.com/mondegor/print-shop-back/internal/factory/filestation/section/pub"
	provideraccountpub "github.com/mondegor/print-shop-back/internal/factory/provideraccounts/section/pub"
)

// RegisterRestRouterPubHandlers - регистрирует в указанном роутере обработчики секции PublicAPI.
func RegisterRestRouterPubHandlers(ctx context.Context, router mrserver.HttpRouter, opts app.Options) error {
	section := NewAppSectionPublicAPI(ctx, opts)
	prepareHandler := mrfactory.WithMiddlewareCheckAccess(ctx, section, opts.AccessControl)
	router.HandlerFunc(http.MethodGet, section.BuildPath("/"), mrresp.HandlerGetStatusOkAsJSON())

	for _, createFunc := range getPublicAPIControllers(ctx, opts) {
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
