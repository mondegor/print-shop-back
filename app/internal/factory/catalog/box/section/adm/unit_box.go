package adm

import (
	"github.com/mondegor/go-storage/mrpostgres/builder"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/catalog/box/section/adm/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/catalog/box/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/catalog/box/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/catalog/box/section/adm/usecase"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

func initBoxController(
	logger mrlog.Logger,
	eventEmitter mrevent.Emitter,
	useCaseErrorWrapper mrerr.UseCaseErrorWrapper,
	dbConnManager mrstorage.DBConnManager,
	requestExtendParser *validate.ExtendParser,
	responseSender mrserver.ResponseSender,
	pageSizeMax uint64,
) (mrserver.HttpController, error) {
	entityMeta, err := mrsql.ParseEntity(logger, entity.Box{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewBoxPostgres(
		dbConnManager,
		builder.NewSQL(
			builder.WithSQLSetMetaEntity(entityMeta.MetaUpdate()),
			builder.WithSQLOrderByDefaultSort(entityMeta.MetaOrderBy().DefaultSort()),
			builder.WithSQLLimitMaxSize(pageSizeMax),
		),
	)

	useCase := usecase.NewBox(storage, eventEmitter, useCaseErrorWrapper)

	controller := httpv1.NewBox(
		requestExtendParser,
		responseSender,
		useCase,
		entityMeta.MetaOrderBy(),
	)

	return controller, nil
}
