package adm

import (
	"github.com/mondegor/go-components/wire/mrordering"
	"github.com/mondegor/go-storage/mrpostgres/builder"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/module"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/usecase"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/shared/validate"
	"github.com/mondegor/print-shop-back/pkg/controls/api"
)

func initFormElementController(
	logger mrlog.Logger,
	eventEmitter mrevent.Emitter,
	dbConnManager mrstorage.DBConnManager,
	storageSubmitForm *repository.SubmitFormPostgres,
	storageFormElement *repository.FormElementPostgres,
	requestParser *validate.Parser,
	responseSender mrserver.ResponseSender,
	elementTemplateAPI api.ElementTemplateHeader,
) (mrserver.HttpController, error) {
	useCase := usecase.NewFormElement(
		storageFormElement,
		storageSubmitForm,
		elementTemplateAPI,
		mrordering.InitServiceMover(
			dbConnManager,
			eventEmitter,
			mrsql.DBTableInfo{
				Name:       module.DBTableNameSubmitFormElements,
				PrimaryKey: "form_id",
			},
		),
		eventEmitter,
		logger,
	)

	controller := httpv1.NewFormElement(
		requestParser,
		responseSender,
		useCase,
	)

	return controller, nil
}

func initFormElementStorage(
	logger mrlog.Logger,
	dbConnManager mrstorage.DBConnManager,
) (*repository.FormElementPostgres, error) {
	entityMeta, err := mrsql.ParseEntity(logger, entity.FormElement{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewFormElementPostgres(
		dbConnManager,
		builder.NewSQL(
			builder.WithSQLSetMetaEntity(entityMeta.MetaUpdate()),
		),
	)

	return storage, nil
}
