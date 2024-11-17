package submitform

import (
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
		EventEmitter        mrsender.EventEmitter
		UseCaseErrorWrapper mrcore.UseCaseErrorWrapper
		DBConnManager       mrstorage.DBConnManager
		Locker              mrlock.Locker
		RequestParsers      RequestParsers
		ResponseSender      mrserver.FileResponseSender

		ElementTemplateAPI api.ElementTemplateHeader

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
