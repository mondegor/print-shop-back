package filestation

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrpath"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

type (
	// Options - comment struct.
	Options struct {
		Logger              mrlog.Logger
		UseCaseErrorWrapper mrerr.UseCaseErrorWrapper
		RequestParser       *mrparser.String
		ResponseSender      mrserver.FileResponseSender

		UnitImageProxy UnitImageProxyOptions
	}

	// UnitImageProxyOptions - comment struct.
	UnitImageProxyOptions struct {
		FileAPI  mrstorage.FileProviderAPI
		BasePath mrpath.PathBuilder
	}
)
