package pub

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/pub/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/pub/usecase"
	"github.com/mondegor/print-shop-back/pkg/transport/validate"
)

func initPrintFormatController(
	dbConnManager mrstorage.DBConnManager,
	requestParser *validate.Parser,
	responseSender mrserver.ResponseSender,
) (mrserver.HttpController, error) {
	storage := repository.NewPrintFormatPostgres(
		dbConnManager,
	)

	useCase := usecase.NewPrintFormat(storage)

	controller := httpv1.NewPrintFormat(
		requestParser,
		responseSender,
		useCase,
	)

	return controller, nil
}
