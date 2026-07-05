package pub

import (
	"github.com/mondegor/go-sysmess/mrstorage"
	"github.com/mondegor/go-webcore/mrcore/initing"
	"github.com/mondegor/go-webcore/mrserver"

	"print-shop-back/internal/controls/submitform/module"
	"print-shop-back/internal/controls/submitform/shared/validate"
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
					return initSubmitFormController(
						dbConnManager,
						requestParser,
						responseSender,
					)
				},
			},
		},
	}
}
