package factory

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresponse"
	"github.com/mondegor/go-webcore/mrtool"
)

type (
	Options struct {
		Logger         mrcore.Logger
		ServiceHelper  *mrtool.ServiceHelper
		RequestParser  mrserver.RequestParserPath
		ResponseSender *mrresponse.Sender

		UnitImageProxy *UnitImageProxyOptions
	}

	UnitImageProxyOptions struct {
		FileAPI mrstorage.FileProviderAPI
		BaseURL string
	}
)
