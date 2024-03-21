package factory

import (
	"context"
	"print-shop-back/internal"
	factory_catalog_box_adm "print-shop-back/internal/modules/catalog/box/factory/admin-api"
	factory_catalog_box_pub "print-shop-back/internal/modules/catalog/box/factory/public-api"
	factory_catalog_laminate_adm "print-shop-back/internal/modules/catalog/laminate/factory/admin-api"
	factory_catalog_laminate_pub "print-shop-back/internal/modules/catalog/laminate/factory/public-api"
	factory_catalog_paper_adm "print-shop-back/internal/modules/catalog/paper/factory/admin-api"
	factory_catalog_paper_pub "print-shop-back/internal/modules/catalog/paper/factory/public-api"
	factory_controls_elementtemplate_adm "print-shop-back/internal/modules/controls/element-template/factory/admin-api"
	factory_controls_submitform_adm "print-shop-back/internal/modules/controls/submit-form/factory/admin-api"
	factory_dictionaries_laminatetype_adm "print-shop-back/internal/modules/dictionaries/laminate-type/factory/admin-api"
	factory_dictionaries_laminatetype_pub "print-shop-back/internal/modules/dictionaries/laminate-type/factory/public-api"
	factory_dictionaries_papercolor_adm "print-shop-back/internal/modules/dictionaries/paper-color/factory/admin-api"
	factory_dictionaries_papercolor_pub "print-shop-back/internal/modules/dictionaries/paper-color/factory/public-api"
	factory_dictionaries_paperfacture_adm "print-shop-back/internal/modules/dictionaries/paper-facture/factory/admin-api"
	factory_dictionaries_paperfacture_pub "print-shop-back/internal/modules/dictionaries/paper-facture/factory/public-api"
	factory_dictionaries_printformat_adm "print-shop-back/internal/modules/dictionaries/print-format/factory/admin-api"
	factory_filestation_pub "print-shop-back/internal/modules/file-station/factory/public-api"
	factory_provider_accounts_adm "print-shop-back/internal/modules/provider-accounts/factory/admin-api"
	factory_provider_accounts_prov "print-shop-back/internal/modules/provider-accounts/factory/providers-api"
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

	if err = RegisterSystemHandlers(ctx, opts.Cfg, router, sectionAdminAPI); err != nil {
		return nil, err
	}

	err = registerAdminAPIControllers(
		ctx,
		opts,
		func(list []mrserver.HttpController, err error) error {
			if err == nil {
				router.Register(
					mrfactory.WithMiddlewareCheckAccess(ctx, list, sectionAdminAPI, opts.AccessControl)...,
				)
			}

			return err
		},
	)

	if err != nil {
		return nil, err
	}

	// section: providers-api
	sectionProvidersAPI := NewAppSectionProvidersAPI(ctx, opts)

	if err = RegisterSystemHandlers(ctx, opts.Cfg, router, sectionProvidersAPI); err != nil {
		return nil, err
	}

	err = registerProvidersAPIControllers(
		ctx,
		opts,
		func(list []mrserver.HttpController, err error) error {
			if err == nil {
				router.Register(
					mrfactory.WithMiddlewareCheckAccess(ctx, list, sectionProvidersAPI, opts.AccessControl)...,
				)
			}

			return err
		},
	)

	if err != nil {
		return nil, err
	}

	// section: public
	sectionPublicAPI := NewAppSectionPublicAPI(ctx, opts)

	if err = RegisterSystemHandlers(ctx, opts.Cfg, router, sectionPublicAPI); err != nil {
		return nil, err
	}

	err = registerPublicAPIControllers(
		ctx,
		opts,
		func(list []mrserver.HttpController, err error) error {
			if err == nil {
				router.Register(
					mrfactory.WithMiddlewareCheckAccess(ctx, list, sectionPublicAPI, opts.AccessControl)...,
				)
			}

			return err
		},
	)

	if err != nil {
		return nil, err
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

func registerAdminAPIControllers(ctx context.Context, opts app.Options, registerFunc func([]mrserver.HttpController, error) error) error {
	if err := registerFunc(factory_catalog_box_adm.CreateModule(ctx, opts.CatalogBoxModule)); err != nil {
		return err
	}

	if err := registerFunc(factory_catalog_laminate_adm.CreateModule(ctx, opts.CatalogLaminateModule)); err != nil {
		return err
	}

	if err := registerFunc(factory_catalog_paper_adm.CreateModule(ctx, opts.CatalogPaperModule)); err != nil {
		return err
	}

	if err := registerFunc(factory_controls_elementtemplate_adm.CreateModule(ctx, opts.ControlsElementTemplateModule)); err != nil {
		return err
	}

	if err := registerFunc(factory_controls_submitform_adm.CreateModule(ctx, opts.ControlsSubmitFormModule)); err != nil {
		return err
	}

	if err := registerFunc(factory_dictionaries_laminatetype_adm.CreateModule(ctx, opts.DictionariesLaminateTypeModule)); err != nil {
		return err
	}

	if err := registerFunc(factory_dictionaries_papercolor_adm.CreateModule(ctx, opts.DictionariesPaperColorModule)); err != nil {
		return err
	}

	if err := registerFunc(factory_dictionaries_paperfacture_adm.CreateModule(ctx, opts.DictionariesPaperFactureModule)); err != nil {
		return err
	}

	if err := registerFunc(factory_dictionaries_printformat_adm.CreateModule(ctx, opts.DictionariesPrintFormatModule)); err != nil {
		return err
	}

	if err := registerFunc(factory_provider_accounts_adm.CreateModule(ctx, opts.ProviderAccountsModule)); err != nil {
		return err
	}

	return nil
}

func registerProvidersAPIControllers(ctx context.Context, opts app.Options, registerFunc func([]mrserver.HttpController, error) error) error {
	if err := registerFunc(factory_provider_accounts_prov.CreateModule(ctx, opts.ProviderAccountsModule)); err != nil {
		return err
	}

	return nil
}

func registerPublicAPIControllers(ctx context.Context, opts app.Options, registerFunc func([]mrserver.HttpController, error) error) error {
	if err := registerFunc(factory_catalog_box_pub.CreateModule(ctx, opts.CatalogBoxModule)); err != nil {
		return err
	}

	if err := registerFunc(factory_catalog_laminate_pub.CreateModule(ctx, opts.CatalogLaminateModule)); err != nil {
		return err
	}

	if err := registerFunc(factory_catalog_paper_pub.CreateModule(ctx, opts.CatalogPaperModule)); err != nil {
		return err
	}

	if err := registerFunc(factory_dictionaries_laminatetype_pub.CreateModule(ctx, opts.DictionariesLaminateTypeModule)); err != nil {
		return err
	}

	if err := registerFunc(factory_dictionaries_papercolor_pub.CreateModule(ctx, opts.DictionariesPaperColorModule)); err != nil {
		return err
	}

	if err := registerFunc(factory_dictionaries_paperfacture_pub.CreateModule(ctx, opts.DictionariesPaperFactureModule)); err != nil {
		return err
	}

	if err := registerFunc(factory_filestation_pub.CreateModule(ctx, opts.FileStationModule)); err != nil {
		return err
	}

	if err := registerFunc(factory_provider_account_pub.CreateModule(ctx, opts.ProviderAccountsModule)); err != nil {
		return err
	}

	return nil
}
