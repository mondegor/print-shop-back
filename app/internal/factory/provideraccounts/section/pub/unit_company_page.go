package pub

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrpath"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/pub/repository"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/pub/usecase"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/shared/validate"
)

func initCompanyPageController(
	useCaseErrorWrapper mrerr.UseCaseErrorWrapper,
	dbConnManager mrstorage.DBConnManager,
	requestModuleParser *validate.Parser,
	responseSender mrserver.ResponseSender,
	logoURLBuilder mrpath.PathBuilder,
) (mrserver.HttpController, error) {
	storage := repository.NewCompanyPagePostgres(dbConnManager)
	useCase := usecase.NewCompanyPage(storage, useCaseErrorWrapper)
	controller := httpv1.NewCompanyPage(
		requestModuleParser,
		responseSender,
		useCase,
		logoURLBuilder,
	)

	return controller, nil
}
