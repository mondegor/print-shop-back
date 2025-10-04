package elementtemplate

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/shared/validate"
)

type (
	// Options - comment struct.
	Options struct {
		Logger               mrlog.Logger
		EventEmitter         mrevent.Emitter
		UsecaseErrorWrapper  mrerr.UseCaseErrorWrapper
		FileUserErrorWrapper mrerr.UserErrorWrapper
		DBConnManager        mrstorage.DBConnManager
		RequestParsers       RequestParsers
		ResponseSender       mrserver.FileResponseSender
		PageSizeMax          uint64
		PageSizeDefault      uint64
	}

	// RequestParsers - comment struct.
	RequestParsers struct {
		// Parser       *pkgvalidate.Parser
		// ExtendParser *pkgvalidate.ExtendParser
		ModuleParser *validate.Parser
	}
)
