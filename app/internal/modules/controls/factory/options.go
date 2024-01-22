package factory

import (
	view_shared "print-shop-back/internal/modules/controls/controller/http_v1/shared/view"
	usecase_api "print-shop-back/internal/modules/controls/usecase/api"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver/mrresponse"
	"github.com/mondegor/go-webcore/mrtool"
)

type (
	Options struct {
		Logger          mrcore.Logger
		EventBox        mrcore.EventBox
		ServiceHelper   *mrtool.ServiceHelper
		PostgresAdapter *mrpostgres.ConnAdapter
		RequestParser   *view_shared.Parser
		ResponseSender  *mrresponse.Sender

		ElementTemplateAPI usecase_api.ElementTemplateAPI
		OrdererAPI         mrorderer.API

		UnitElementTemplate *UnitElementTemplateOptions
	}

	UnitElementTemplateOptions struct {
		Dictionary *mrlang.MultiLangDictionary
	}
)
