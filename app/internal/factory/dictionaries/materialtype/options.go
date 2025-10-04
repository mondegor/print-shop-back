package materialtype

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/pkg/validate"
)

type (
	// Options - comment struct.
	Options struct {
		Logger              mrlog.Logger
		EventEmitter        mrevent.Emitter
		UsecaseErrorWrapper mrerr.UseCaseErrorWrapper
		DBConnManager       mrstorage.DBConnManager
		RequestParsers      RequestParsers
		ResponseSender      mrserver.ResponseSender

		PageSizeMax     uint64
		PageSizeDefault uint64
	}

	// RequestParsers - comment struct.
	RequestParsers struct {
		Parser       *validate.Parser
		ExtendParser *validate.ExtendParser
	}
)
