package adm

import (
	"github.com/mondegor/go-storage/mrpostgres/builder"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-webcore/mrserver"

	"print-shop-back/internal/adapter/log"
	"print-shop-back/internal/dictionaries/papercolor/section/adm/controller/httpv1"
	"print-shop-back/internal/dictionaries/papercolor/section/adm/entity"
	"print-shop-back/internal/dictionaries/papercolor/section/adm/repository"
	"print-shop-back/internal/dictionaries/papercolor/section/adm/usecase"
	"print-shop-back/pkg/transport/validate"
)

func initPaperColorController(
	logger log.Logger,
	eventEmitter mrevent.Emitter,
	dbConnManager mrstorage.DBConnManager,
	requestExtendParser *validate.ExtendParser,
	responseSender mrserver.ResponseSender,
	pageSizeMax int,
) (mrserver.HttpController, error) {
	entityMeta, err := mrsql.ParseEntity(logger, entity.PaperColor{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewPaperColorPostgres(
		dbConnManager,
		builder.NewSQL(
			builder.WithSQLOrderByDefaultSort(entityMeta.MetaOrderBy().DefaultSort()),
			builder.WithSQLLimitMaxSize(pageSizeMax),
		),
	)

	useCase := usecase.NewPaperColor(storage, eventEmitter)

	controller := httpv1.NewPaperColor(
		requestExtendParser,
		responseSender,
		useCase,
		entityMeta.MetaOrderBy(),
	)

	return controller, nil
}
