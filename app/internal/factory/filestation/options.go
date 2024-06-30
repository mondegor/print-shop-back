package filestation

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

type (
	// Options - comment struct.
	Options struct {
		UsecaseHelper  mrcore.UsecaseErrorWrapper
		RequestParser  *mrparser.String
		ResponseSender mrserver.FileResponseSender

		UnitImageProxy UnitImageProxyOptions
	}

	// UnitImageProxyOptions - comment struct.
	UnitImageProxyOptions struct {
		FileAPI mrstorage.FileProviderAPI
		BaseURL string
	}
)
