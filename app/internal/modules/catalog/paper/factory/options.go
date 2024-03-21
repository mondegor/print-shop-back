package factory

import (
	view_shared "print-shop-back/internal/modules/catalog/paper/controller/http_v1/shared/view"
	"print-shop-back/pkg/modules/dictionaries"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrserver/mrresponse"
)

type (
	Options struct {
		EventEmitter    mrsender.EventEmitter
		UsecaseHelper   *mrcore.UsecaseHelper
		PostgresAdapter *mrpostgres.ConnAdapter
		RequestParser   *view_shared.Parser
		ResponseSender  *mrresponse.Sender

		PaperColorAPI   dictionaries.PaperColorAPI
		PaperFactureAPI dictionaries.PaperFactureAPI

		UnitPaper UnitPaperOptions

		PageSizeMax     uint64
		PageSizeDefault uint64
	}

	UnitPaperOptions struct {
		Dictionary *mrlang.MultiLangDictionary
	}
)
