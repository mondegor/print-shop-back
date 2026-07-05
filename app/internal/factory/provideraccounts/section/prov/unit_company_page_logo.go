package prov

import (
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlock"
	"github.com/mondegor/go-sysmess/mrstorage"
	"github.com/mondegor/go-webcore/mrserver"

	"print-shop-back/internal/adapter/log"
	"print-shop-back/internal/provideraccounts/section/prov/controller/httpv1"
	"print-shop-back/internal/provideraccounts/section/prov/repository"
	"print-shop-back/internal/provideraccounts/section/prov/usecase"
	"print-shop-back/internal/provideraccounts/shared/validate"
)

func initCompanyPageLogoController(
	logger log.Logger,
	eventEmitter mrevent.Emitter,
	dbConnManager mrstorage.DBConnManager,
	locker mrlock.Locker,
	requestModuleParser *validate.Parser,
	responseSender mrserver.ResponseSender,
	logoFileAPI mrstorage.FileProviderAPI,
) (mrserver.HttpController, error) {
	storage := repository.NewCompanyPageLogoPostgres(dbConnManager)

	useCase := usecase.NewCompanyPageLogo(
		storage,
		logoFileAPI,
		locker,
		eventEmitter,
		logger,
	)

	controller := httpv1.NewCompanyPageLogo(
		requestModuleParser,
		responseSender,
		useCase,
	)

	return controller, nil
}
