package app

import (
	"net/http"

	"github.com/mondegor/go-components/mrmailer"
	"github.com/mondegor/go-components/mrnotifier"
	"github.com/mondegor/go-components/mrsettings"
	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrredis"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlocale"
	"github.com/mondegor/go-sysmess/mrlock"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtrace"
	"github.com/mondegor/go-sysmess/mrwire"
	"github.com/mondegor/go-webcore/mraccess"
	"github.com/mondegor/go-webcore/mraccess/role/filestorage"
	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrpath"
	"github.com/mondegor/go-webcore/mrrun"
	"github.com/mondegor/go-webcore/mrsentry"
	"github.com/mondegor/go-webcore/mrserver/mrhttp"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
	"github.com/mondegor/go-webcore/mrserver/mrresp"
	"github.com/mondegor/go-webcore/mrworker/process/collect"
	"github.com/mondegor/go-webcore/mrworker/process/consume"
	"github.com/mondegor/go-webcore/mrworker/process/schedule"

	"github.com/mondegor/print-shop-back/config"
	"github.com/mondegor/print-shop-back/internal/factory/auth"
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
		Cfg             config.Config
		Logger          mrlog.Logger
		Tracer          mrtrace.Tracer
		TraceManager    mrtrace.ContextManager
		OpenedResources *mrwire.CloseManager

		InternalRouter        *http.ServeMux
		Sentry                *mrsentry.Adapter
		Prometheus            *mrinit.Prometheus
		EventEmitter          mrevent.Emitter
		ErrorHandler          mrerr.ErrorHandler
		UsecaseErrorWrapper   mrerr.UseCaseErrorWrapper
		StorageErrorWrapper   mrerr.ErrorWrapper
		FileUserErrorWrapper  mrerr.UserErrorWrapper
		ImageUserErrorWrapper mrerr.UserErrorWrapper
		AppHealth             *mrrun.AppHealth

		PostgresConnManager          *mrpostgres.ConnManager
		PostgresNotificationService  *mrpostgres.ProcessWaitForNotification
		PostgresNotificationChannels mrpostgres.ReceiverChannels
		RedisAdapter                 *mrredis.ConnAdapter
		FileProviderPool             *mrstorage.FileProviderPool
		Locker                       mrlock.Locker
		LocalePool                   *mrlocale.Pool
		RequestParsers               RequestParsers
		ResponseSenders              ResponseSenders
		PermsProvider                *filestorage.PermsProvider
		RealmKindRights              mraccess.RightsGetter
		ImageURLBuilder              mrpath.PathBuilder

		// API section
		DictionariesMaterialTypeAPI api.MaterialTypeAvailability
		DictionariesPaperColorAPI   api.PaperColorAvailability
		DictionariesPaperFactureAPI api.PaperFactureAvailability
		DictionariesPrintFormatAPI  api.PrintFormatAvailability
		MailerAPI                   mrmailer.MessageProducer
		NotifierAPI                 mrnotifier.NoticeProducer
		SettingsGetterAPI           mrsettings.DefaultValueGetter
		SettingsSetterAPI           mrsettings.Setter

		// Modules section
		AuthModule                     auth.Options
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

		// Services and Servers section
		UserStatRequestCollectorService *collect.MessageCollector
		MailProcessorService            *consume.MessageProcessor
		NoticeProcessorService          *consume.MessageProcessor
		HttpServer                      *mrhttp.Adapter
		HttpInternalServer              *mrhttp.Adapter
		TaskSchedulerServices           []*schedule.TaskScheduler // можно добавлять начиная с формирования API
	}

	// RequestParsers - comment struct.
	RequestParsers struct {
		// Bool       *mrparser.Bool
		// DateTime   *mrparser.DateTime
		Int64      *mrparser.Int64
		ItemStatus *mrparser.ItemStatus
		Uint64     *mrparser.Uint64
		ListSorter *mrparser.ListSorter
		ListPager  *mrparser.ListPager
		String     *mrparser.String
		UUID       *mrparser.UUID
		Validator  *mrparser.Validator
		Locale     *mrparser.Locale
		ClientIP   *mrparser.ClientIP
		User       *mrparser.User
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
