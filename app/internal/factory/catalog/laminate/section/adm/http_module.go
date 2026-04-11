package adm

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrcore/initing"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/catalog/laminate/module"
	"github.com/mondegor/print-shop-back/pkg/dictionaries/api"
	"github.com/mondegor/print-shop-back/pkg/transport/validate"
)

// InitHttpModule - создаются все компоненты модуля и возвращаются к нему контролеры.
func InitHttpModule(
	logger mrlog.Logger,
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
