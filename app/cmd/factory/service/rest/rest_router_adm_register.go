package rest

import (
	"net/http"

	"github.com/mondegor/go-webcore/mraccess"
	"github.com/mondegor/go-webcore/mraccess/section"
	"github.com/mondegor/go-webcore/mrcore/mrinit"
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
func RegisterRestRouterAdmHandlers(
	router mrserver.HttpRouter,
	opts app.Options,
	sect *section.RoutingSection,
	memberProvider mraccess.MemberProvider,
) error {
	router.HandlerFunc(http.MethodGet, sect.BuildPath("/"), mrresp.HandlerGetStatusOkAsJSON(opts.Logger))
	prepareHandler := mrinit.WithMiddlewareCheckAccess(opts.Logger, sect, memberProvider, opts.RealmKindRights, opts.PermsProvider)

	for _, createFunc := range getAdminAPIControllers(opts) {
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

func getAdminAPIControllers(opts app.Options) []func() (list []mrserver.HttpController, err error) {
	return []func() (list []mrserver.HttpController, err error){
		func() ([]mrserver.HttpController, error) {
			return catalogboxadm.CreateModule(opts.CatalogBoxModule)
		},
		func() ([]mrserver.HttpController, error) {
			return cataloglaminateadm.CreateModule(opts.CatalogLaminateModule)
		},
		func() ([]mrserver.HttpController, error) {
			return catalogpaperadm.CreateModule(opts.CatalogPaperModule)
		},
		func() ([]mrserver.HttpController, error) {
			return controlselementtemplateadm.CreateModule(opts.ControlsElementTemplateModule)
		},
		func() ([]mrserver.HttpController, error) {
			return controlssubmitformadm.CreateModule(opts.ControlsSubmitFormModule)
		},
		func() ([]mrserver.HttpController, error) {
			return dictionariesmaterialtypeadm.CreateModule(opts.DictionariesMaterialTypeModule)
		},
		func() ([]mrserver.HttpController, error) {
			return dictionariespapercoloradm.CreateModule(opts.DictionariesPaperColorModule)
		},
		func() ([]mrserver.HttpController, error) {
			return dictionariespaperfactureadm.CreateModule(opts.DictionariesPaperFactureModule)
		},
		func() ([]mrserver.HttpController, error) {
			return dictionariesprintformatadm.CreateModule(opts.DictionariesPrintFormatModule)
		},
		func() ([]mrserver.HttpController, error) {
			return provideraccountsadm.CreateModule(opts.ProviderAccountsModule)
		},
	}
}
