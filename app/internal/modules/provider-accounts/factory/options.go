package factory

import (
	view_shared "print-shop-back/internal/modules/provider-accounts/controller/http_v1/shared/view"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlock"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
	"github.com/mondegor/go-webcore/mrserver/mrresponse"
)

type (
	Options struct {
		EventEmitter    mrsender.EventEmitter
		UsecaseHelper   *mrcore.UsecaseHelper
		PostgresAdapter *mrpostgres.ConnAdapter
		Locker          mrlock.Locker
		RequestParsers  RequestParsers
		ResponseSender  *mrresponse.Sender

		UnitCompanyPage UnitCompanyPageOptions

		PageSizeMax     uint64
		PageSizeDefault uint64
	}

	UnitCompanyPageOptions struct {
		LogoFileAPI    mrstorage.FileProviderAPI
		LogoURLBuilder mrlib.BuilderPath
	}

	RequestParsers struct {
		String *mrparser.String
		Image  *mrparser.Image
		Parser *view_shared.Parser
	}
)
