package pub

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/pub/repository"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/pub/usecase"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/shared/validate"
)

func initSubmitFormController(
	dbConnManager mrstorage.DBConnManager,
	requestParser *validate.Parser,
	responseSender mrserver.ResponseSender,
) (mrserver.HttpController, error) {
	storage := repository.NewSubmitFormPostgres(
		dbConnManager,
	)

	useCase := usecase.NewSubmitForm(storage)

	controller := httpv1.NewSubmitForm(
		requestParser,
		responseSender,
		useCase,
	)

	return controller, nil
}
