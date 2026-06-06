package provideraccounts

import (
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlock"
	"github.com/mondegor/go-sysmess/mrpath"
	"github.com/mondegor/go-sysmess/mrstorage"
	"github.com/mondegor/go-webcore/mrserver"

	"print-shop-back/internal/adapter/log"
	"print-shop-back/internal/provideraccounts/shared/validate"
)

type (
	// Options - comment struct.
	Options struct {
		Logger         log.Logger
		EventEmitter   mrevent.Emitter
		DBConnManager  mrstorage.DBConnManager
		Locker         mrlock.Locker
		RequestParsers RequestParsers
		ResponseSender mrserver.ResponseSender

		UnitCompanyPage UnitCompanyPageOptions

		PageSizeMax     int
		PageSizeDefault int
	}

	// UnitCompanyPageOptions - comment struct.
	UnitCompanyPageOptions struct {
		LogoFileAPI    mrstorage.FileProviderAPI
		LogoURLBuilder mrpath.Builder
	}

	// RequestParsers - comment struct.
	RequestParsers struct {
		// Parser       *pkgvalidate.Parser
		// ExtendParser *pkgvalidate.ExtendParser
		ModuleParser *validate.Parser
	}
)
