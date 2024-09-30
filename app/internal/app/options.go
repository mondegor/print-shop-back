package app

import (
	"context"
	"net/http"

	"github.com/mondegor/go-components/mrsettings"
	"github.com/mondegor/go-components/mrsort"
	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrredis"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrlock"
	"github.com/mondegor/go-webcore/mrpath"
	"github.com/mondegor/go-webcore/mrperms"
	"github.com/mondegor/go-webcore/mrrun"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrsentry"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
	"github.com/mondegor/go-webcore/mrserver/mrresp"
	"github.com/mondegor/go-webcore/mrworker"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/mondegor/print-shop-back/config"
	calculationsalgo "github.com/mondegor/print-shop-back/internal/factory/calculations/algo"
	calculationsquery "github.com/mondegor/print-shop-back/internal/factory/calculations/queryhistory"
	catalogbox "github.com/mondegor/print-shop-back/internal/factory/catalog/box"
	cataloglaminate "github.com/mondegor/print-shop-back/internal/factory/catalog/laminate"
	catalogpaper "github.com/mondegor/print-shop-back/internal/factory/catalog/paper"
	controlselementtemplate "github.com/mondegor/print-shop-back/internal/factory/controls/elementtemplate"
	controlssubmitform "github.com/mondegor/print-shop-back/internal/factory/controls/submitform"
	dictionariesmaterialtype "github.com/mondegor/print-shop-back/internal/factory/dictionaries/materialtype"
	dictionariespapercolor "github.com/mondegor/print-shop-back/internal/factory/dictionaries/papercolor"
	dictionariespaperfacture "github.com/mondegor/print-shop-back/internal/factory/dictionaries/paperfacture"
	dictionariesprintformat "github.com/mondegor/print-shop-back/internal/factory/dictionaries/printformat"
	"github.com/mondegor/print-shop-back/internal/factory/filestation"
	"github.com/mondegor/print-shop-back/internal/factory/provideraccounts"
	"github.com/mondegor/print-shop-back/pkg/dictionaries/api"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

type (
	// Options - comment struct.
	Options struct {
		Cfg                 config.Config
		AppHealth           *mrrun.AppHealth
		ErrorHandler        mrcore.ErrorHandler
		EventEmitter        mrsender.EventEmitter
		ErrorManager        *mrinit.ErrorManager
		UseCaseErrorWrapper mrcore.UseCaseErrorWrapper

		InternalRouter *http.ServeMux
		Sentry         *mrsentry.Adapter
		Prometheus     *prometheus.Registry

		PostgresConnManager *mrpostgres.ConnManager
		RedisAdapter        *mrredis.ConnAdapter
		FileProviderPool    *mrstorage.FileProviderPool
		Locker              mrlock.Locker
		Translator          *mrlang.Translator
		RequestParsers      RequestParsers
		ResponseSenders     ResponseSenders
		AccessControl       *mrperms.RoleAccessControl
		ImageURLBuilder     mrpath.PathBuilder

		// API section
		DictionariesMaterialTypeAPI api.MaterialTypeAvailability
		DictionariesPaperColorAPI   api.PaperColorAvailability
		DictionariesPaperFactureAPI api.PaperFactureAvailability
		DictionariesPrintFormatAPI  api.PrintFormatAvailability
		OrdererAPI                  mrsort.Orderer
		SettingsGetterAPI           mrsettings.DefaultValueGetter
		SettingsSetterAPI           mrsettings.Setter

		// Modules section
		CalculationsAlgoModule         calculationsalgo.Options
		CalculationsQueryHistoryModule calculationsquery.Options
		CatalogBoxModule               catalogbox.Options
		CatalogLaminateModule          cataloglaminate.Options
		CatalogPaperModule             catalogpaper.Options
		ControlsElementTemplateModule  controlselementtemplate.Options
		ControlsSubmitFormModule       controlssubmitform.Options
		DictionariesMaterialTypeModule dictionariesmaterialtype.Options
		DictionariesPaperColorModule   dictionariespapercolor.Options
		DictionariesPaperFactureModule dictionariespaperfacture.Options
		DictionariesPrintFormatModule  dictionariesprintformat.Options
		FileStationModule              filestation.Options
		ProviderAccountsModule         provideraccounts.Options

		SchedulerTasks  []mrworker.Task
		OpenedResources []func(ctx context.Context)
	}

	// RequestParsers - comment struct.
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

		Parser       *validate.Parser
		ExtendParser *validate.ExtendParser
	}

	// ResponseSenders - comment struct.
	ResponseSenders struct {
		Sender     *mrresp.Sender
		FileSender *mrresp.FileSender
	}
)
