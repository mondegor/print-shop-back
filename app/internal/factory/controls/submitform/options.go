package submitform

import (
	"github.com/mondegor/go-components/mrsort"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlock"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/shared/validate"
	"github.com/mondegor/print-shop-back/pkg/controls/api"
)

type (
	// Options - comment struct.
	Options struct {
		EventEmitter   mrsender.EventEmitter
		UsecaseHelper  mrcore.UsecaseErrorWrapper
		DBConnManager  mrstorage.DBConnManager
		Locker         mrlock.Locker
		RequestParsers RequestParsers
		ResponseSender mrserver.FileResponseSender

		ElementTemplateAPI api.ElementTemplateHeader
		OrdererAPI         mrsort.Orderer

		PageSizeMax     uint64
		PageSizeDefault uint64
	}

	// RequestParsers - comment struct.
	RequestParsers struct {
		// Parser       *pkgvalidate.Parser
		// ExtendParser *pkgvalidate.ExtendParser
		ModuleParser *validate.Parser
	}
)
