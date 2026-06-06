package usr

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-webcore/mrcore/initing"
	"github.com/mondegor/go-webcore/mrserver"

	"print-shop-back/internal/adapter/log"
	"print-shop-back/internal/warehousing/module"
	"print-shop-back/pkg/transport/validate"
)

// InitHttpModule - создаются все компоненты модуля и возвращаются к нему контролеры.
func InitHttpModule(
	logger log.Logger,
	eventEmitter mrevent.Emitter,
	dbConnManager mrstorage.DBConnManager,
	requestExtendParser *validate.ExtendParser,
	responseSender mrserver.ResponseSender,
) initing.HttpModule {
	return initing.HttpModule{
		Caption:    module.Name,
		Permission: module.Permission,
		Controllers: []initing.HttpController{
			{
				Create: func() (mrserver.HttpController, error) {
					return initContainerController(
						logger,
						eventEmitter,
						dbConnManager,
						requestExtendParser,
						responseSender,
					)
				},
			},
			{
				Create: func() (mrserver.HttpController, error) {
					return initStockController(
						logger,
						eventEmitter,
						dbConnManager,
						requestExtendParser,
						responseSender,
					)
				},
			},
			{
				Create: func() (mrserver.HttpController, error) {
					return initStoreController(
						dbConnManager,
						requestExtendParser,
						responseSender,
					)
				},
			},
		},
	}
}
