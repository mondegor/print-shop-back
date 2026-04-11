package prov

import (
	"github.com/mondegor/go-storage/mrlock"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrpath"
	"github.com/mondegor/go-webcore/mrcore/initing"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/module"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/shared/validate"
)

// InitHttpModule - создаются все компоненты модуля и возвращаются к нему контролеры.
func InitHttpModule(
	logger mrlog.Logger,
	eventEmitter mrevent.Emitter,
	dbConnManager mrstorage.DBConnManager,
	locker mrlock.Locker,
	requestModuleParser *validate.Parser,
	responseSender mrserver.ResponseSender,
	logoFileAPIFunc func() (mrstorage.FileProviderAPI, error),
	logoURLBuilder mrpath.Builder,
) initing.HttpModule {
	return initing.HttpModule{
		Caption:    module.Name,
		Permission: module.Permission,
		Controllers: []initing.HttpController{
			{
				Create: func() (mrserver.HttpController, error) {
					return initCompanyPageController(
						eventEmitter,
						dbConnManager,
						requestModuleParser,
						responseSender,
						logoURLBuilder,
					)
				},
			},
			{
				Create: func() (mrserver.HttpController, error) {
					logoFileAPI, err := logoFileAPIFunc()
					if err != nil {
						return nil, err
					}

					return initCompanyPageLogoController(
						logger,
						eventEmitter,
						dbConnManager,
						locker,
						requestModuleParser,
						responseSender,
						logoFileAPI,
					)
				},
			},
		},
	}
}
