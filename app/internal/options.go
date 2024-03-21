package app

import (
	"context"
	"print-shop-back/config"
	factory_catalog_box "print-shop-back/internal/modules/catalog/box/factory"
	factory_catalog_laminate "print-shop-back/internal/modules/catalog/laminate/factory"
	factory_catalog_paper "print-shop-back/internal/modules/catalog/paper/factory"
	factory_controls_elementtemplate "print-shop-back/internal/modules/controls/element-template/factory"
	factory_controls_submitform "print-shop-back/internal/modules/controls/submit-form/factory"
	factory_dictionaries_laminatetype "print-shop-back/internal/modules/dictionaries/laminate-type/factory"
	factory_dictionaries_papercolor "print-shop-back/internal/modules/dictionaries/paper-color/factory"
	factory_dictionaries_paperfacture "print-shop-back/internal/modules/dictionaries/paper-facture/factory"
	factory_dictionaries_printformat "print-shop-back/internal/modules/dictionaries/print-format/factory"
	factory_filestation "print-shop-back/internal/modules/file-station/factory"
	factory_provider_accounts "print-shop-back/internal/modules/provider-accounts/factory"
	"print-shop-back/pkg/modules/dictionaries"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrredis"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlock"
	"github.com/mondegor/go-webcore/mrperms"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
	"github.com/mondegor/go-webcore/mrserver/mrresponse"
)

type (
	Options struct {
		Cfg              config.Config
		EventEmitter     mrsender.EventEmitter
		UsecaseHelper    *mrcore.UsecaseHelper
		PostgresAdapter  *mrpostgres.ConnAdapter
		RedisAdapter     *mrredis.ConnAdapter
		FileProviderPool *mrstorage.FileProviderPool
		Locker           mrlock.Locker
		Translator       *mrlang.Translator
		RequestParsers   RequestParsers
		ResponseSender   *mrresponse.Sender
		AccessControl    mrperms.AccessControl
		ImageURLBuilder  mrlib.BuilderPath

		// API section
		DictionariesLaminateTypeAPI dictionaries.LaminateTypeAPI
		DictionariesPaperColorAPI   dictionaries.PaperColorAPI
		DictionariesPaperFactureAPI dictionaries.PaperFactureAPI
		DictionariesPrintFormatAPI  dictionaries.PrintFormatAPI
		OrdererAPI                  mrorderer.API

		// Modules section
		CatalogBoxModule               factory_catalog_box.Options
		CatalogLaminateModule          factory_catalog_laminate.Options
		CatalogPaperModule             factory_catalog_paper.Options
		ControlsElementTemplateModule  factory_controls_elementtemplate.Options
		ControlsSubmitFormModule       factory_controls_submitform.Options
		DictionariesLaminateTypeModule factory_dictionaries_laminatetype.Options
		DictionariesPaperColorModule   factory_dictionaries_papercolor.Options
		DictionariesPaperFactureModule factory_dictionaries_paperfacture.Options
		DictionariesPrintFormatModule  factory_dictionaries_printformat.Options
		FileStationModule              factory_filestation.Options
		ProviderAccountsModule         factory_provider_accounts.Options

		OpenedResources []func(ctx context.Context)
	}

	RequestParsers struct {
		// Bool       *mrparser.Bool
		// DateTime   *mrparser.DateTime
		Int64      *mrparser.Int64
		ItemStatus *mrparser.ItemStatus
		KeyInt32   *mrparser.KeyInt32
		ListSorter *mrparser.ListSorter
		ListPager  *mrparser.ListPager
		String     *mrparser.String
		UUID       *mrparser.UUID
		Validator  *mrparser.Validator
		FileJson   *mrparser.File
		ImageLogo  *mrparser.Image
	}
)
