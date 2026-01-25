package adm

import (
	"github.com/mondegor/go-storage/mrpostgres/builder"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/adm/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/adm/usecase"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

func initMaterialTypeController(
	logger mrlog.Logger,
	eventEmitter mrevent.Emitter,
	dbConnManager mrstorage.DBConnManager,
	requestExtendParser *validate.ExtendParser,
	responseSender mrserver.ResponseSender,
	pageSizeMax uint64,
) (mrserver.HttpController, error) {
	entityMeta, err := mrsql.ParseEntity(logger, entity.MaterialType{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewMaterialTypePostgres(
		dbConnManager,
		builder.NewSQL(
			builder.WithSQLOrderByDefaultSort(entityMeta.MetaOrderBy().DefaultSort()),
			builder.WithSQLLimitMaxSize(pageSizeMax),
		),
	)

	useCase := usecase.NewMaterialType(storage, eventEmitter)

	controller := httpv1.NewMaterialType(
		requestExtendParser,
		responseSender,
		useCase,
		entityMeta.MetaOrderBy(),
	)

	return controller, nil
}
