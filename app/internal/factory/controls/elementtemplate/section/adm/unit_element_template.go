package adm

import (
	"github.com/mondegor/go-storage/mrpostgres/builder"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/section/adm/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/section/adm/usecase"
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/shared/validate"
)

func initElementTemplateController(
	logger mrlog.Logger,
	eventEmitter mrevent.Emitter,
	useCaseErrorWrapper mrerr.UseCaseErrorWrapper,
	fileUserErrorWrapper mrerr.UserErrorWrapper,
	dbConnManager mrstorage.DBConnManager,
	requestParser *validate.Parser,
	responseSender mrserver.FileResponseSender,
	pageSizeMax uint64,
) (mrserver.HttpController, error) {
	entityMeta, err := mrsql.ParseEntity(logger, entity.ElementTemplate{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewElementTemplatePostgres(
		dbConnManager,
		builder.NewSQL(
			builder.WithSQLSetMetaEntity(entityMeta.MetaUpdate()),
			builder.WithSQLOrderByDefaultSort(entityMeta.MetaOrderBy().DefaultSort()),
			builder.WithSQLLimitMaxSize(pageSizeMax),
		),
	)

	useCase := usecase.NewElementTemplate(storage, eventEmitter, useCaseErrorWrapper)

	controller := httpv1.NewElementTemplate(
		requestParser,
		responseSender,
		useCase,
		entityMeta.MetaOrderBy(),
		fileUserErrorWrapper,
	)

	return controller, nil
}
