package pub

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrpath"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"

	"github.com/mondegor/print-shop-back/internal/filestation/module"
	"github.com/mondegor/print-shop-back/internal/initing"
)

// InitHttpModule - создаются все компоненты модуля и возвращаются к нему контролеры.
func InitHttpModule(
	useCaseErrorWrapper mrerr.UseCaseErrorWrapper,
	requestParser *mrparser.String,
	responseFileSender mrserver.FileResponseSender,
	fileAPIFunc func() (mrstorage.FileProviderAPI, mrpath.PathBuilder, error),
) initing.HttpModule {
	return initing.HttpModule{
		Name:       module.Name,
		Permission: module.Permission,
		Controllers: []initing.HttpController{
			{
				Name:       module.UnitImageProxyName,
				Permission: module.UnitImageProxyPermission,
				Create: func() (mrserver.HttpController, error) {
					fileAPI, basePath, err := fileAPIFunc()
					if err != nil {
						return nil, err
					}

					return initImageProxyController(
						useCaseErrorWrapper,
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
