package adm

import (
	"github.com/mondegor/go-storage/mrpostgres/builder"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-webcore/mrserver"

	"print-shop-back/internal/adapter/log"
	"print-shop-back/internal/dictionaries/materialtype/section/adm/controller/httpv1"
	"print-shop-back/internal/dictionaries/materialtype/section/adm/entity"
	"print-shop-back/internal/dictionaries/materialtype/section/adm/repository"
	"print-shop-back/internal/dictionaries/materialtype/section/adm/usecase"
	"print-shop-back/pkg/transport/validate"
)

func initMaterialTypeController(
	logger log.Logger,
	eventEmitter mrevent.Emitter,
	dbConnManager mrstorage.DBConnManager,
	requestExtendParser *validate.ExtendParser,
	responseSender mrserver.ResponseSender,
	pageSizeMax int,
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
