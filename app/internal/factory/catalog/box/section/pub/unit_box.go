package pub

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/catalog/box/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/catalog/box/section/pub/repository"
	"github.com/mondegor/print-shop-back/internal/catalog/box/section/pub/usecase"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

func initBoxController(
	dbConnManager mrstorage.DBConnManager,
	requestParser *validate.Parser,
	responseSender mrserver.ResponseSender,
) (mrserver.HttpController, error) {
	storage := repository.NewBoxPostgres(
		dbConnManager,
	)

	useCase := usecase.NewBox(storage)

	controller := httpv1.NewBox(
		requestParser,
		responseSender,
		useCase,
	)

	return controller, nil
}
