package pub

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/pub/repository"
	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/pub/usecase"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

func initLaminateController(
	useCaseErrorWrapper mrerr.UseCaseErrorWrapper,
	dbConnManager mrstorage.DBConnManager,
	requestParser *validate.Parser,
	responseSender mrserver.ResponseSender,
) (mrserver.HttpController, error) {
	storage := repository.NewLaminatePostgres(
		dbConnManager,
	)

	useCase := usecase.NewLaminate(storage, useCaseErrorWrapper)

	controller := httpv1.NewLaminate(
		requestParser,
		responseSender,
		useCase,
	)

	return controller, nil
}
