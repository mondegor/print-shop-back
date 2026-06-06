package pub

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-webcore/mrserver"

	"print-shop-back/internal/calculations/queryhistory/section/pub/controller/httpv1"
	"print-shop-back/internal/calculations/queryhistory/section/pub/repository"
	"print-shop-back/internal/calculations/queryhistory/section/pub/usecase"
	"print-shop-back/pkg/transport/validate"
)

func initQueryHistoryController(
	eventEmitter mrevent.Emitter,
	dbConnManager mrstorage.DBConnManager,
	requestExtendParser *validate.ExtendParser,
	responseSender mrserver.ResponseSender,
) (mrserver.HttpController, error) {
	storage := repository.NewQueryHistoryPostgres(
		dbConnManager,
	)

	useCase := usecase.NewQueryHistory(storage, eventEmitter)

	controller := httpv1.NewQueryHistory(
		requestExtendParser,
		responseSender,
		useCase,
	)

	return controller, nil
}
