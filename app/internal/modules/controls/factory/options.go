package factory

import (
	view_shared "print-shop-back/internal/modules/controls/controller/http_v1/shared/view"
	usecase_api "print-shop-back/internal/modules/controls/usecase/api"

	"github.com/mondegor/go-components/mrorderer"
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

		ElementTemplateAPI usecase_api.ElementTemplateAPI
		OrdererAPI         mrorderer.API

		UnitElementTemplate UnitElementTemplateOptions

		PageSizeMax     uint64
		PageSizeDefault uint64
	}

	UnitElementTemplateOptions struct {
		Dictionary *mrlang.MultiLangDictionary
	}
)
