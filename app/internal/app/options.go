package app

import (
	"net/http"

	"github.com/mondegor/go-components/mrmailer"
	"github.com/mondegor/go-components/mrnotifier"
	"github.com/mondegor/go-components/mrsettings"
	"github.com/mondegor/go-storage/mrlock"
	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrredis"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlocale"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrpath"
	"github.com/mondegor/go-sysmess/mrtrace"
	"github.com/mondegor/go-sysmess/util/xio"
	"github.com/mondegor/go-webcore/mraccess"
	"github.com/mondegor/go-webcore/mrclient/sentry"
	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrrun"
	"github.com/mondegor/go-webcore/mrserver/httpserver"
	"github.com/mondegor/go-webcore/mrserver/mrresp"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
	"github.com/mondegor/go-webcore/mrworker/process/collect"
	"github.com/mondegor/go-webcore/mrworker/process/consume"

	"github.com/mondegor/print-shop-back/config"
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
		OpenedResources *xio.CloseManager

		InternalRouter *http.ServeMux
		Sentry         *sentry.Adapter
		Prometheus     *mrinit.Prometheus
		EventEmitter   mrevent.Emitter
		ErrorHandler   errors.Handler
		AppHealth      *mrrun.AppHealth

		PostgresConnManager          *mrpostgres.ConnManager
		PostgresNotificationService  *mrpostgres.ProcessWaitForNotification
		PostgresNotificationChannels mrpostgres.ReceiverChannels
		RedisAdapter                 *mrredis.ConnAdapter
		FileProviderPool             *mrstorage.FileProviderPool
		Locker                       mrlock.Locker
		LocalePool                   *mrlocale.Pool
		RequestParsers               RequestParsers
		ResponseSenders              ResponseSenders
		PermsProvider                mraccess.RightsSource
		RealmUserProviders           map[string]mraccess.UserProvider
		ImageURLBuilder              mrpath.Builder

		// API section
		DictionariesMaterialTypeAPI api.MaterialTypeAvailability
		DictionariesPaperColorAPI   api.PaperColorAvailability
		DictionariesPaperFactureAPI api.PaperFactureAvailability
		DictionariesPrintFormatAPI  api.PrintFormatAvailability
		MailerAPI                   mrmailer.MessageProducer
		NotifierAPI                 mrnotifier.NoteProducer
		SettingsGetterAPI           mrsettings.MustGetter
		SettingsSetterAPI           mrsettings.Setter

		// Services and Servers section
		UserStatRequestCollectorService *collect.MessageCollector
		MailProcessorService            *consume.MessageProcessor
		NoticeProcessorService          *consume.MessageProcessor
		HttpServer                      *httpserver.Adapter
		HttpInternalServer              *httpserver.Adapter
		TaskSchedulerServices           []mrrun.Process // можно добавлять начиная с формирования API
	}

	// RequestParsers - comment struct.
	RequestParsers struct {
		// Bool       *parser.Bool
		// DateTime   *parser.DateTime
		Int64      *parser.Int64
		ItemStatus *parser.ItemStatus
		Uint64     *parser.Uint64
		ListSorter *parser.ListSorter
		ListPager  *parser.ListPager
		String     *parser.String
		UUID       *parser.UUID
		Validator  *parser.Validator
		Locale     *parser.Locale
		ClientIP   *parser.ClientIP
		User       *parser.User
		FileJson   *parser.File
		ImageLogo  *parser.Image

		Parser       *validate.Parser
		ExtendParser *validate.ExtendParser
	}

	// ResponseSenders - comment struct.
	ResponseSenders struct {
		Sender     *mrresp.Sender
		FileSender *mrresp.FileSender
	}
)
