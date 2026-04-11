package pub

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrpath"
	"github.com/mondegor/go-webcore/mrcore/initing"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/request/parser"

	"github.com/mondegor/print-shop-back/internal/filestation/module"
)

// InitHttpModule - создаются все компоненты модуля и возвращаются к нему контролеры.
func InitHttpModule(
	requestParser *parser.String,
	responseFileSender mrserver.FileResponseSender,
	fileAPIFunc func() (mrstorage.FileProviderAPI, mrpath.Builder, error),
) initing.HttpModule {
	return initing.HttpModule{
		Caption:    module.Name,
		Permission: module.Permission,
		Controllers: []initing.HttpController{
			{
				Caption:    module.UnitImageProxyName,
				Permission: module.UnitImageProxyPermission,
				Create: func() (mrserver.HttpController, error) {
					fileAPI, basePath, err := fileAPIFunc()
					if err != nil {
						return nil, err
					}

					return initImageProxyController(
						requestParser,
						responseFileSender,
						fileAPI,
						basePath,
					)
				},
			},
		},
	}
}
