package adm

import (
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlock"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/module"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/shared/validate"
	"github.com/mondegor/print-shop-back/internal/initing"
	"github.com/mondegor/print-shop-back/pkg/controls/api"
)

// InitHttpModule - создаются все компоненты модуля и возвращаются к нему контролеры.
func InitHttpModule(
	logger mrlog.Logger,
	eventEmitter mrevent.Emitter,
	useCaseErrorWrapper mrerr.UseCaseErrorWrapper,
	storageErrorWrapper mrerr.ErrorWrapper,
	dbConnManager mrstorage.DBConnManager,
	locker mrlock.Locker,
	requestModuleParser *validate.Parser,
	responseSender mrserver.ResponseSender,
	responseFileSender mrserver.FileResponseSender,
	elementTemplateAPI api.ElementTemplateHeader,
	pageSizeMax uint64,
) initing.HttpModule {
	// переменные, которые должны быть инициализированы в InitSharedComponents модуля
	var (
		storageSubmitForm    *repository.SubmitFormPostgres
		entityMetaSubmitForm *mrsql.EntityMetaOrderBy
		storageFormElement   *repository.FormElementPostgres
	)

	return initing.HttpModule{
		Name:       module.Name,
		Permission: module.Permission,
		InitSharedComponents: func() (err error) {
			storageSubmitForm, entityMetaSubmitForm, err = initSubmitFormStorage(logger, dbConnManager, pageSizeMax)
			if err != nil {
				return err
			}

			storageFormElement, err = initFormElementStorage(logger, dbConnManager)

			return err
		},
		Controllers: []initing.HttpController{
			{
				Create: func() (mrserver.HttpController, error) {
					return initSubmitFormController(
						eventEmitter,
						useCaseErrorWrapper,
						dbConnManager,
						storageSubmitForm,
						storageFormElement,
						entityMetaSubmitForm,
						locker,
						requestModuleParser,
						responseFileSender,
					)
				},
			},
			{
				Name:       module.UnitFormElementName,
				Permission: module.UnitFormElementPermission,
				Create: func() (mrserver.HttpController, error) {
					return initFormElementController(
						logger,
						eventEmitter,
						useCaseErrorWrapper,
						storageErrorWrapper,
						dbConnManager,
						storageSubmitForm,
						storageFormElement,
						requestModuleParser,
						responseSender,
						elementTemplateAPI,
					)
				},
			},
		},
	}
}
