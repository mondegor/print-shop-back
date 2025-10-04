package rest

import (
	"net/http"

	"github.com/mondegor/go-webcore/mraccess"
	"github.com/mondegor/go-webcore/mraccess/section"
	"github.com/mondegor/go-webcore/mrcore/mrinit"
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
func RegisterRestRouterPubHandlers(router mrserver.HttpRouter, opts app.Options, sect *section.RoutingSection, memberProvider mraccess.MemberProvider) error {
	router.HandlerFunc(http.MethodGet, sect.BuildPath("/"), mrresp.HandlerGetStatusOkAsJSON(opts.Logger))
	prepareHandler := mrinit.WithMiddlewareCheckAccess(opts.Logger, sect, memberProvider, opts.RealmKindRights, opts.PermsProvider)

	for _, createFunc := range getPublicAPIControllers(opts) {
		list, err := createFunc()
		if err != nil {
			return err
		}

		router.Register(
			mrinit.PrepareEachController(list, prepareHandler)...,
		)
	}

	return nil
}

func getPublicAPIControllers(opts app.Options) []func() (list []mrserver.HttpController, err error) {
	return []func() (list []mrserver.HttpController, err error){
		func() ([]mrserver.HttpController, error) {
			return calculationsquerypub.CreateModule(opts.CalculationsQueryHistoryModule)
		},
		func() ([]mrserver.HttpController, error) {
			return calculationsalgopub.CreateModule(opts.CalculationsAlgoModule)
		},
		func() ([]mrserver.HttpController, error) {
			return catalogboxpub.CreateModule(opts.CatalogBoxModule)
		},
		func() ([]mrserver.HttpController, error) {
			return cataloglaminatepub.CreateModule(opts.CatalogLaminateModule)
		},
		func() ([]mrserver.HttpController, error) {
			return catalogpaperpub.CreateModule(opts.CatalogPaperModule)
		},
		func() ([]mrserver.HttpController, error) {
			return controlssubmitformpub.CreateModule(opts.ControlsSubmitFormModule)
		},
		func() ([]mrserver.HttpController, error) {
			return dictionariesmaterialtypepub.CreateModule(opts.DictionariesMaterialTypeModule)
		},
		func() ([]mrserver.HttpController, error) {
			return dictionariespapercolorpub.CreateModule(opts.DictionariesPaperColorModule)
		},
		func() ([]mrserver.HttpController, error) {
			return dictionariespaperfacturepub.CreateModule(opts.DictionariesPaperFactureModule)
		},
		func() ([]mrserver.HttpController, error) {
			return dictionariesprintformatpub.CreateModule(opts.DictionariesPrintFormatModule)
		},
		func() ([]mrserver.HttpController, error) {
			return filestationpub.CreateModule(opts.FileStationModule)
		},
		func() ([]mrserver.HttpController, error) {
			return provideraccountpub.CreateModule(opts.ProviderAccountsModule)
		},
	}
}
