package factory

import (
	view_shared "print-shop-back/internal/modules/controls/submit-form/controller/http_v1/shared/view"
	"print-shop-back/pkg/modules/controls"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlock"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrserver/mrresponse"
)

type (
	Options struct {
		EventEmitter    mrsender.EventEmitter
		UsecaseHelper   *mrcore.UsecaseHelper
		PostgresAdapter *mrpostgres.ConnAdapter
		Locker          mrlock.Locker
		RequestParser   *view_shared.Parser
		ResponseSender  *mrresponse.Sender

		ElementTemplateAPI controls.ElementTemplateAPI
		OrdererAPI         mrorderer.API

		PageSizeMax     uint64
		PageSizeDefault uint64
	}
)
