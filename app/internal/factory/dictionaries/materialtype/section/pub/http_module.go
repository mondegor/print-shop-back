package pub

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore/initing"
	"github.com/mondegor/go-webcore/mrserver"

	"print-shop-back/internal/dictionaries/materialtype/module"
	"print-shop-back/pkg/transport/validate"
)

// InitHttpModule - создаются все компоненты модуля и возвращаются к нему контролеры.
func InitHttpModule(
	dbConnManager mrstorage.DBConnManager,
	requestParser *validate.Parser,
	responseSender mrserver.ResponseSender,
) initing.HttpModule {
	return initing.HttpModule{
		Caption:    module.Name,
		Permission: module.Permission,
		Controllers: []initing.HttpController{
			{
				Create: func() (mrserver.HttpController, error) {
					return initMaterialTypeController(
						dbConnManager,
						requestParser,
						responseSender,
					)
				},
			},
		},
	}
}
