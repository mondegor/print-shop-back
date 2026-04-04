package pub

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/pub/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/pub/usecase"
	"github.com/mondegor/print-shop-back/pkg/transport/validate"
)

func initPaperColorController(
	dbConnManager mrstorage.DBConnManager,
	requestParser *validate.Parser,
	responseSender mrserver.ResponseSender,
) (mrserver.HttpController, error) {
	storage := repository.NewPaperColorPostgres(
		dbConnManager,
	)

	useCase := usecase.NewPaperColor(storage)

	controller := httpv1.NewPaperColor(
		requestParser,
		responseSender,
		useCase,
	)

	return controller, nil
}
