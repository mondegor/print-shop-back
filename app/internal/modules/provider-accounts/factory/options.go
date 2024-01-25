package factory

import (
	view_shared "print-shop-back/internal/modules/provider-accounts/controller/http_v1/shared/view"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
	"github.com/mondegor/go-webcore/mrserver/mrresponse"
	"github.com/mondegor/go-webcore/mrtool"
)

type (
	Options struct {
		Logger          mrcore.Logger
		EventBox        mrcore.EventBox
		ServiceHelper   *mrtool.ServiceHelper
		PostgresAdapter *mrpostgres.ConnAdapter
		Locker          mrcore.Locker
		RequestParsers  *RequestParsers
		ResponseSender  *mrresponse.Sender

		UnitCompanyPage *UnitCompanyPageOptions
	}

	UnitCompanyPageOptions struct {
		LogoFileAPI    mrstorage.FileProviderAPI
		LogoURLBuilder mrcore.BuilderPath
	}

	RequestParsers struct {
		String *mrparser.String
		Image  *mrparser.Image
		Parser *view_shared.Parser
	}
)
