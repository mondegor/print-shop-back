package adm

import (
	"github.com/mondegor/go-storage/mrpostgres/builder"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-webcore/mrserver"

	"print-shop-back/internal/adapter/log"
	"print-shop-back/internal/dictionaries/paperfacture/section/adm/controller/httpv1"
	"print-shop-back/internal/dictionaries/paperfacture/section/adm/entity"
	"print-shop-back/internal/dictionaries/paperfacture/section/adm/repository"
	"print-shop-back/internal/dictionaries/paperfacture/section/adm/usecase"
	"print-shop-back/pkg/transport/validate"
)

func initPaperFactureController(
	logger log.Logger,
	eventEmitter mrevent.Emitter,
	dbConnManager mrstorage.DBConnManager,
	requestExtendParser *validate.ExtendParser,
	responseSender mrserver.ResponseSender,
	pageSizeMax int,
) (mrserver.HttpController, error) {
	entityMeta, err := mrsql.ParseEntity(logger, entity.PaperFacture{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewPaperFacturePostgres(
		dbConnManager,
		builder.NewSQL(
			builder.WithSQLOrderByDefaultSort(entityMeta.MetaOrderBy().DefaultSort()),
			builder.WithSQLLimitMaxSize(pageSizeMax),
		),
	)

	useCase := usecase.NewPaperFacture(storage, eventEmitter)

	controller := httpv1.NewPaperFacture(
		requestExtendParser,
		responseSender,
		useCase,
		entityMeta.MetaOrderBy(),
	)

	return controller, nil
}
