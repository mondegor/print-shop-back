package adm

import (
	"github.com/mondegor/go-storage/mrpostgres/builder"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlock"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/usecase"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/shared/validate"
)

func initSubmitFormController(
	eventEmitter mrevent.Emitter,
	useCaseErrorWrapper mrerr.UseCaseErrorWrapper,
	dbConnManager mrstorage.DBConnManager,
	storageSubmitForm *repository.SubmitFormPostgres,
	storageFormElement *repository.FormElementPostgres,
	entityMetaSubmitForm *mrsql.EntityMetaOrderBy,
	locker mrlock.Locker,
	requestParser *validate.Parser,
	responseSender mrserver.FileResponseSender,
) (mrserver.HttpController, error) {
	storageFormVersion := repository.NewFormVersionPostgres(
		dbConnManager,
	)

	useCase := usecase.NewSubmitForm(
		storageSubmitForm,
		storageFormElement,
		storageFormVersion,
		eventEmitter,
		useCaseErrorWrapper,
	)

	useCaseVersion := usecase.NewFormVersion(
		storageFormVersion,
		useCase,
		usecase.NewFormCompilerJson(),
		locker,
		eventEmitter,
		useCaseErrorWrapper,
	)

	controller := httpv1.NewSubmitForm(
		requestParser,
		responseSender,
		useCase,
		useCaseVersion,
		entityMetaSubmitForm,
	)

	return controller, nil
}

func initSubmitFormStorage(
	logger mrlog.Logger,
	dbConnManager mrstorage.DBConnManager,
	pageSizeMax uint64,
) (*repository.SubmitFormPostgres, *mrsql.EntityMetaOrderBy, error) {
	entityMeta, err := mrsql.ParseEntity(logger, entity.SubmitForm{})
	if err != nil {
		return nil, nil, err
	}

	storage := repository.NewSubmitFormPostgres(
		dbConnManager,
		builder.NewSQL(
			builder.WithSQLSetMetaEntity(entityMeta.MetaUpdate()),
			builder.WithSQLOrderByDefaultSort(entityMeta.MetaOrderBy().DefaultSort()),
			builder.WithSQLLimitMaxSize(pageSizeMax),
		),
	)

	return storage, entityMeta.MetaOrderBy(), nil
}
