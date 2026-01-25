package prov

import (
	"github.com/mondegor/go-storage/mrlock"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/prov/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/prov/repository"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/prov/usecase"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/shared/validate"
)

func initCompanyPageLogoController(
	logger mrlog.Logger,
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
