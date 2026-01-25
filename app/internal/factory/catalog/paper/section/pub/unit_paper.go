package pub

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/pub/repository"
	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/pub/usecase"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

func initPaperController(
	dbConnManager mrstorage.DBConnManager,
	requestParser *validate.Parser,
	responseSender mrserver.ResponseSender,
) (mrserver.HttpController, error) {
	storage := repository.NewPaperPostgres(
		dbConnManager,
	)

	useCase := usecase.NewPaper(storage)

	controller := httpv1.NewPaper(
		requestParser,
		responseSender,
		useCase,
	)

	return controller, nil
}
