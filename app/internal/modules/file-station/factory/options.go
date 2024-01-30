package factory

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
	"github.com/mondegor/go-webcore/mrserver/mrresponse"
)

type (
	Options struct {
		UsecaseHelper  *mrcore.UsecaseHelper
		RequestParser  *mrparser.String
		ResponseSender *mrresponse.Sender

		UnitImageProxy UnitImageProxyOptions
	}

	UnitImageProxyOptions struct {
		FileAPI mrstorage.FileProviderAPI
		BaseURL string
	}
)
