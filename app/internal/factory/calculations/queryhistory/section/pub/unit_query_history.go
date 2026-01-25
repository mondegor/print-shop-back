package pub

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/queryhistory/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/calculations/queryhistory/section/pub/repository"
	"github.com/mondegor/print-shop-back/internal/calculations/queryhistory/section/pub/usecase"
	"github.com/mondegor/print-shop-back/pkg/validate"
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
