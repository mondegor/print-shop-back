package elementtemplate

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/shared/validate"
)

type (
	// Options - comment struct.
	Options struct {
		EventEmitter   mrsender.EventEmitter
		UsecaseHelper  mrcore.UsecaseErrorWrapper
		DBConnManager  mrstorage.DBConnManager
		RequestParsers RequestParsers
		ResponseSender mrserver.FileResponseSender

		UnitElementTemplate UnitElementTemplateOptions

		PageSizeMax     uint64
		PageSizeDefault uint64
	}

	// UnitElementTemplateOptions - comment struct.
	UnitElementTemplateOptions struct {
		Dictionary *mrlang.MultiLangDictionary
	}

	// RequestParsers - comment struct.
	RequestParsers struct {
		// Parser       *pkgvalidate.Parser
		// ExtendParser *pkgvalidate.ExtendParser
		ModuleParser *validate.Parser
	}
)
