package factory

import (
	view_shared "print-shop-back/internal/modules/catalog/laminate/controller/http_v1/shared/view"
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

		LaminateTypeAPI dictionaries.LaminateTypeAPI

		UnitLaminate UnitLaminateOptions

		PageSizeMax     uint64
		PageSizeDefault uint64
	}

	UnitLaminateOptions struct {
		Dictionary *mrlang.MultiLangDictionary
	}
)
