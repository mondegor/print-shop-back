package adm

import (
	"github.com/mondegor/go-core/mrevent"
	"github.com/mondegor/go-core/mrstorage"
	"github.com/mondegor/go-webcore/mrcore/initing"
	"github.com/mondegor/go-webcore/mrserver"

	"print-shop-back/internal/adapter/log"
	"print-shop-back/internal/catalog/laminate/module"
	"print-shop-back/pkg/dictionaries/api"
	"print-shop-back/pkg/transport/validate"
)

// InitHttpModule - создаются все компоненты модуля и возвращаются к нему контролеры.
func InitHttpModule(
	logger log.Logger,
	eventEmitter mrevent.Emitter,
	dbConnManager mrstorage.DBConnManager,
	requestExtendParser *validate.ExtendParser,
	responseSender mrserver.ResponseSender,
	materialTypeAPI api.MaterialTypeAvailability,
	pageSizeMax int,
) initing.HttpModule {
	return initing.HttpModule{
		Caption:    module.Name,
		Permission: module.Permission,
		Controllers: []initing.HttpController{
			{
				Create: func() (mrserver.HttpController, error) {
					return initLaminateController(
						logger,
						eventEmitter,
						dbConnManager,
						requestExtendParser,
						responseSender,
						materialTypeAPI,
						pageSizeMax,
					)
				},
			},
		},
	}
}
