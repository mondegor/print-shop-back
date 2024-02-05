package factory

import (
	"context"
	"print-shop-back/internal"
	factory_catalog_box_adm "print-shop-back/internal/modules/catalog/box/factory/admin-api"
	factory_catalog_laminate_adm "print-shop-back/internal/modules/catalog/laminate/factory/admin-api"
	factory_catalog_paper_adm "print-shop-back/internal/modules/catalog/paper/factory/admin-api"
	factory_controls_adm "print-shop-back/internal/modules/controls/factory/admin-api"
	factory_dictionaries_laminatetype_adm "print-shop-back/internal/modules/dictionaries/laminate-type/factory/admin-api"
	factory_dictionaries_papercolor_adm "print-shop-back/internal/modules/dictionaries/paper-color/factory/admin-api"
	factory_dictionaries_paperfacture_adm "print-shop-back/internal/modules/dictionaries/paper-facture/factory/admin-api"
	factory_dictionaries_printformat_adm "print-shop-back/internal/modules/dictionaries/print-format/factory/admin-api"
	factory_filestation_pub "print-shop-back/internal/modules/file-station/factory/public-api"
	factory_provider_accounts_adm "print-shop-back/internal/modules/provider-accounts/factory/admin-api"
	factory_provider_accounts_pacc "print-shop-back/internal/modules/provider-accounts/factory/provider-account-api"
	factory_provider_account_pub "print-shop-back/internal/modules/provider-accounts/factory/public-api"
	"time"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
)

const (
	restServerCaption = "RestServer"
)

func NewRestServer(ctx context.Context, opts app.Options) (*mrserver.ServerAdapter, error) {
	mrlog.Ctx(ctx).Info().Msgf("Create and init '%s'", restServerCaption)

	router, err := NewRestRouter(ctx, opts.Cfg, opts.Translator)

	if err != nil {
		return nil, err
	}

	// section: admin-api
	sectionAdminAPI := NewAppSectionAdminAPI(ctx, opts)

	if err := RegisterSystemHandlers(ctx, opts.Cfg, router, sectionAdminAPI); err != nil {
		return nil, err
	}

	if controllers, err := factory_catalog_box_adm.CreateModule(ctx, opts.CatalogBoxModule); err != nil {
		return nil, err
	} else {
		router.Register(
			mrfactory.WithMiddlewareCheckAccess(ctx, controllers, sectionAdminAPI, opts.AccessControl)...,
		)
	}

	if controllers, err := factory_catalog_laminate_adm.CreateModule(ctx, opts.CatalogLaminateModule); err != nil {
		return nil, err
	} else {
		router.Register(
			mrfactory.WithMiddlewareCheckAccess(ctx, controllers, sectionAdminAPI, opts.AccessControl)...,
		)
	}

	if controllers, err := factory_catalog_paper_adm.CreateModule(ctx, opts.CatalogPaperModule); err != nil {
		return nil, err
	} else {
		router.Register(
			mrfactory.WithMiddlewareCheckAccess(ctx, controllers, sectionAdminAPI, opts.AccessControl)...,
		)
	}

	if controllers, err := factory_controls_adm.CreateModule(ctx, opts.ControlsModule); err != nil {
		return nil, err
	} else {
		router.Register(
			mrfactory.WithMiddlewareCheckAccess(ctx, controllers, sectionAdminAPI, opts.AccessControl)...,
		)
	}

	if controllers, err := factory_dictionaries_laminatetype_adm.CreateModule(ctx, opts.DictionariesLaminateTypeModule); err != nil {
		return nil, err
	} else {
		router.Register(
			mrfactory.WithMiddlewareCheckAccess(ctx, controllers, sectionAdminAPI, opts.AccessControl)...,
		)
	}

	if controllers, err := factory_dictionaries_papercolor_adm.CreateModule(ctx, opts.DictionariesPaperColorModule); err != nil {
		return nil, err
	} else {
		router.Register(
			mrfactory.WithMiddlewareCheckAccess(ctx, controllers, sectionAdminAPI, opts.AccessControl)...,
		)
	}

	if controllers, err := factory_dictionaries_paperfacture_adm.CreateModule(ctx, opts.DictionariesPaperFactureModule); err != nil {
		return nil, err
	} else {
		router.Register(
			mrfactory.WithMiddlewareCheckAccess(ctx, controllers, sectionAdminAPI, opts.AccessControl)...,
		)
	}

	if controllers, err := factory_dictionaries_printformat_adm.CreateModule(ctx, opts.DictionariesPrintFormatModule); err != nil {
		return nil, err
	} else {
		router.Register(
			mrfactory.WithMiddlewareCheckAccess(ctx, controllers, sectionAdminAPI, opts.AccessControl)...,
		)
	}

	if controllers, err := factory_provider_accounts_adm.CreateModule(ctx, opts.ProviderAccountsModule); err != nil {
		return nil, err
	} else {
		router.Register(
			mrfactory.WithMiddlewareCheckAccess(ctx, controllers, sectionAdminAPI, opts.AccessControl)...,
		)
	}

	// section: provider-account-api
	sectionProviderAccountAPI := NewAppSectionProviderAccountAPI(ctx, opts)

	if err := RegisterSystemHandlers(ctx, opts.Cfg, router, sectionProviderAccountAPI); err != nil {
		return nil, err
	}

	if controllers, err := factory_provider_accounts_pacc.CreateModule(ctx, opts.ProviderAccountsModule); err != nil {
		return nil, err
	} else {
		router.Register(
			mrfactory.WithMiddlewareCheckAccess(ctx, controllers, sectionProviderAccountAPI, opts.AccessControl)...,
		)
	}

	// section: public
	sectionPublicAPI := NewAppSectionPublicAPI(ctx, opts)

	if err := RegisterSystemHandlers(ctx, opts.Cfg, router, sectionPublicAPI); err != nil {
		return nil, err
	}

	if controllers, err := factory_filestation_pub.CreateModule(ctx, opts.FileStationModule); err != nil {
		return nil, err
	} else {
		router.Register(
			mrfactory.WithMiddlewareCheckAccess(ctx, controllers, sectionPublicAPI, opts.AccessControl)...,
		)
	}

	if controllers, err := factory_provider_account_pub.CreateModule(ctx, opts.ProviderAccountsModule); err != nil {
		return nil, err
	} else {
		router.Register(
			mrfactory.WithMiddlewareCheckAccess(ctx, controllers, sectionPublicAPI, opts.AccessControl)...,
		)
	}

	srvOpts := opts.Cfg.Servers.RestServer

	return mrserver.NewServerAdapter(
		ctx,
		mrserver.ServerOptions{
			Caption:         restServerCaption,
			Handler:         router,
			ReadTimeout:     srvOpts.ReadTimeout * time.Second,
			WriteTimeout:    srvOpts.WriteTimeout * time.Second,
			ShutdownTimeout: srvOpts.ShutdownTimeout * time.Second,
			Listen: mrserver.ListenOptions{
				AppPath:  opts.Cfg.AppPath,
				Type:     srvOpts.Listen.Type,
				SockName: srvOpts.Listen.SockName,
				BindIP:   srvOpts.Listen.BindIP,
				Port:     srvOpts.Listen.Port,
			},
		},
	), nil
}
