package pub

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/pub/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/pub/usecase"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

func initPaperFactureController(
	useCaseErrorWrapper mrerr.UseCaseErrorWrapper,
	dbConnManager mrstorage.DBConnManager,
	requestParser *validate.Parser,
	responseSender mrserver.ResponseSender,
) (mrserver.HttpController, error) {
	storage := repository.NewPaperFacturePostgres(
		dbConnManager,
	)

	useCase := usecase.NewPaperFacture(storage, useCaseErrorWrapper)

	controller := httpv1.NewPaperFacture(
		requestParser,
		responseSender,
		useCase,
	)

	return controller, nil
}
