package adm

import (
	"github.com/mondegor/go-storage/mrpostgres/builder"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-webcore/mrserver"

	"print-shop-back/internal/adapter/log"
	"print-shop-back/internal/dictionaries/printformat/section/adm/controller/httpv1"
	"print-shop-back/internal/dictionaries/printformat/section/adm/entity"
	"print-shop-back/internal/dictionaries/printformat/section/adm/repository"
	"print-shop-back/internal/dictionaries/printformat/section/adm/usecase"
	"print-shop-back/pkg/transport/validate"
)

func initPrintFormatController(
	logger log.Logger,
	eventEmitter mrevent.Emitter,
	dbConnManager mrstorage.DBConnManager,
	requestExtendParser *validate.ExtendParser,
	responseSender mrserver.ResponseSender,
	pageSizeMax int,
) (mrserver.HttpController, error) {
	entityMeta, err := mrsql.ParseEntity(logger, entity.PrintFormat{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewPrintFormatPostgres(
		dbConnManager,
		builder.NewSQL(
			builder.WithSQLOrderByDefaultSort(entityMeta.MetaOrderBy().DefaultSort()),
			builder.WithSQLLimitMaxSize(pageSizeMax),
		),
	)

	useCase := usecase.NewPrintFormat(storage, eventEmitter)

	controller := httpv1.NewPrintFormat(
		requestExtendParser,
		responseSender,
		useCase,
		entityMeta.MetaOrderBy(),
	)

	return controller, nil
}
