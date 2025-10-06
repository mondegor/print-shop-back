package adm

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/module"
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/shared/validate"
	"github.com/mondegor/print-shop-back/internal/initing"
)

// InitHttpModule - создаются все компоненты модуля и возвращаются к нему контролеры.
func InitHttpModule(
	logger mrlog.Logger,
	eventEmitter mrevent.Emitter,
	useCaseErrorWrapper mrerr.UseCaseErrorWrapper,
	fileUserErrorWrapper mrerr.UserErrorWrapper,
	dbConnManager mrstorage.DBConnManager,
	requestParser *validate.Parser,
	responseSender mrserver.FileResponseSender,
	pageSizeMax uint64,
) initing.HttpModule {
	return initing.HttpModule{
		Name:       module.Name,
		Permission: module.Permission,
		Controllers: []initing.HttpController{
			{
				Create: func() (mrserver.HttpController, error) {
					return initElementTemplateController(
						logger,
						eventEmitter,
						useCaseErrorWrapper,
						fileUserErrorWrapper,
						dbConnManager,
						requestParser,
						responseSender,
						pageSizeMax,
					)
				},
			},
		},
	}
}
