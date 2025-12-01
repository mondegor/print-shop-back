package prov

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrpath"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/prov/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/prov/repository"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/prov/usecase"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/shared/validate"
)

func initCompanyPageController(
	eventEmitter mrevent.Emitter,
	useCaseErrorWrapper mrerr.UseCaseErrorWrapper,
	dbConnManager mrstorage.DBConnManager,
	requestModuleParser *validate.Parser,
	responseSender mrserver.ResponseSender,
	logoURLBuilder mrpath.Builder,
) (mrserver.HttpController, error) {
	storage := repository.NewCompanyPagePostgres(dbConnManager)

	useCase := usecase.NewCompanyPage(
		dbConnManager,
		storage,
		logoURLBuilder,
		eventEmitter,
		useCaseErrorWrapper,
	)

	controller := httpv1.NewCompanyPage(
		requestModuleParser,
		responseSender,
		useCase,
	)

	return controller, nil
}
