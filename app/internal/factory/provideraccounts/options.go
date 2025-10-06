package provideraccounts

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlock"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrpath"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/shared/validate"
)

type (
	// Options - comment struct.
	Options struct {
		Logger                mrlog.Logger
		EventEmitter          mrevent.Emitter
		UseCaseErrorWrapper   mrerr.UseCaseErrorWrapper
		ImageUserErrorWrapper mrerr.UserErrorWrapper
		DBConnManager         mrstorage.DBConnManager
		Locker                mrlock.Locker
		RequestParsers        RequestParsers
		ResponseSender        mrserver.ResponseSender

		UnitCompanyPage UnitCompanyPageOptions

		PageSizeMax     uint64
		PageSizeDefault uint64
	}

	// UnitCompanyPageOptions - comment struct.
	UnitCompanyPageOptions struct {
		LogoFileAPI    mrstorage.FileProviderAPI
		LogoURLBuilder mrpath.PathBuilder
	}

	// RequestParsers - comment struct.
	RequestParsers struct {
		// Parser       *pkgvalidate.Parser
		// ExtendParser *pkgvalidate.ExtendParser
		ModuleParser *validate.Parser
	}
)
