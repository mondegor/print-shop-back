package app

import (
	"net/http"

	"github.com/mondegor/go-components/mrauth/dto"
	"github.com/mondegor/go-components/mrmailer"
	mailerentity "github.com/mondegor/go-components/mrmailer/entity"
	"github.com/mondegor/go-components/mrnotifier"
	notifierentity "github.com/mondegor/go-components/mrnotifier/notifier/entity"
	"github.com/mondegor/go-components/mrsettings"
	"github.com/mondegor/go-storage/mrredis"
	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mraccess"
	"github.com/mondegor/go-sysmess/mraccess/provider/filestorage"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlocale"
	"github.com/mondegor/go-sysmess/mrlock"
	"github.com/mondegor/go-sysmess/mrpath"
	"github.com/mondegor/go-sysmess/mrpostgres"
	"github.com/mondegor/go-sysmess/mrpostgres/listennotify"
	"github.com/mondegor/go-sysmess/mrprocess/collect"
	"github.com/mondegor/go-sysmess/mrprocess/consume"
	"github.com/mondegor/go-sysmess/mrrun"
	"github.com/mondegor/go-sysmess/mrstorage"
	"github.com/mondegor/go-sysmess/util/xio"
	"github.com/mondegor/go-webcore/mrclient/sentry"
	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrserver/httpserver"
	"github.com/mondegor/go-webcore/mrserver/mrresp"
	"github.com/mondegor/go-webcore/mrserver/request/parser"

	"print-shop-back/config"
	"print-shop-back/internal/adapter/log"
	"print-shop-back/internal/adapter/trace"
	"print-shop-back/pkg/dictionaries/api"
	validate2 "print-shop-back/pkg/transport/validate"
)

type (
	// Options - comment struct.
	Options struct {
		Cfg             config.Config
		Logger          log.Logger
		Tracer          trace.Tracer
		TraceManager    trace.ContextManager
		OpenedResources *xio.CloseManager
		DebugFunc       func(value any) string

		MonitoringRouter *http.ServeMux
		Sentry           *sentry.Adapter
		Prometheus       *mrinit.Prometheus
		EventEmitter     mrevent.Emitter
		ErrorHandler     errors.Handler
		AppHealth        *mrrun.AppHealth

		PostgresConnManager         *mrpostgres.ConnManager
		PostgresNotificationService *listennotify.ProcessWaitForNotification
		RedisAdapter                *mrredis.ConnAdapter
		FileProviderPool            *mrstorage.FileProviderPool
		Locker                      mrlock.Locker
		LocalePool                  *mrlocale.Pool
		RequestParsers              RequestParsers
		ResponseSenders             ResponseSenders
		PermsProvider               *filestorage.PermsProvider
		RealmUserProviders          map[string]mraccess.UserProvider
		ImageURLBuilder             mrpath.Builder

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
		UserStatRequestCollectorService *collect.MessageCollector[dto.UserActivityLogMessage]
		MailProcessorService            *consume.MessageProcessor[mailerentity.Message]
		NoticeProcessorService          *consume.MessageProcessor[notifierentity.Note]
		HttpServer                      *httpserver.Adapter
		HttpMonitoringServer            *httpserver.Adapter
		TaskSchedulerServices           []mrrun.Process // можно добавлять начиная с формирования API
	}

	// RequestParsers - comment struct.
	RequestParsers struct {
		// Bool       *parser.Bool
		// DateTime   *parser.DateTime
		Int64      *parser.Int64
		ItemStatus *parser.ItemStatus
		Uint64     *parser.Uint64
		ListCursor *parser.ListCursor
		ListPager  *parser.ListPager
		ListSorter *parser.ListSorter
		String     *parser.String
		UUID       *parser.UUID
		Validator  *parser.Validator
		Locale     *parser.Locale
		ClientIP   *parser.ClientIP
		User       *parser.User
		FileJson   *parser.File
		ImageLogo  *parser.Image

		Parser       *validate2.Parser
		ExtendParser *validate2.ExtendParser
	}

	// ResponseSenders - comment struct.
	ResponseSenders struct {
		Sender     *mrresp.Sender
		FileSender *mrresp.FileSender
	}
)
