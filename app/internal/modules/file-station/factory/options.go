package factory

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrtool"
)

type (
	Options struct {
		Logger        mrcore.Logger
		ServiceHelper *mrtool.ServiceHelper

		UnitImageProxy *UnitImageProxyOptions
	}

	UnitImageProxyOptions struct {
		FileAPI mrstorage.FileProviderAPI
		BaseURL string
	}
)
