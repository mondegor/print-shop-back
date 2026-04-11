package pub

import (
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrcore/initing"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/module"
	"github.com/mondegor/print-shop-back/pkg/transport/validate"
)

// InitHttpModule - создаются все компоненты модуля и возвращаются к нему контролеры.
func InitHttpModule(
	logger mrlog.Logger,
	eventEmitter mrevent.Emitter,
	requestParser *validate.Parser,
	responseSender mrserver.ResponseSender,
) initing.HttpModule {
	return initing.HttpModule{
		Caption:    module.Name,
		Permission: module.Permission,
		Controllers: []initing.HttpController{
			{
				Create: func() (mrserver.HttpController, error) {
					return initBoxPackInBoxController(
						logger,
						eventEmitter,
						requestParser,
						responseSender,
					)
				},
			},
			{
				Create: func() (mrserver.HttpController, error) {
					return initBoxSheetCuttingController(
						eventEmitter,
						requestParser,
						responseSender,
					)
				},
			},
			{
				Create: func() (mrserver.HttpController, error) {
					return initSheetImpositionController(
						logger,
						eventEmitter,
						requestParser,
						responseSender,
					)
				},
			},
			{
				Create: func() (mrserver.HttpController, error) {
					return initSheetInsideOutsideController(
						eventEmitter,
						requestParser,
						responseSender,
					)
				},
			},
			{
				Create: func() (mrserver.HttpController, error) {
					return initSheetPackInStackController(
						eventEmitter,
						requestParser,
						responseSender,
					)
				},
			},
		},
	}
}
