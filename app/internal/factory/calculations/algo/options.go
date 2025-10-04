package algo

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

type (
	// Options - comment struct.
	Options struct {
		Logger         mrlog.Logger
		EventEmitter   mrevent.Emitter
		DBConnManager  mrstorage.DBConnManager
		RequestParsers RequestParsers
		ResponseSender mrserver.ResponseSender
	}

	// RequestParsers - comment struct.
	RequestParsers struct {
		Validator *mrparser.Validator
	}
)
