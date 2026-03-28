package usr

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	repository2 "github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/cmd/repository"
	usecase3 "github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/cmd/usecase"
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/repository"
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/service"
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/usecase"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

func initContainerController(
	logger mrlog.Logger,
	eventEmitter mrevent.Emitter,
	dbConnManager mrstorage.DBConnManager,
	requestExtendParser *validate.ExtendParser,
	responseSender mrserver.ResponseSender,
) (mrserver.HttpController, error) {
	storage := repository.NewContainerPostgres(
		dbConnManager,
	)

	// :TODO: временно, вынести на уровень выше
	storage2 := repository.NewStorePostgres(
		dbConnManager,
	)

	// :TODO: временно, вынести на уровень выше
	storageStock := repository.NewStockPostgres(
		dbConnManager,
	)

	// :TODO: временно, вынести на уровень выше
	usecase2 := usecase3.NewRefreshStoreContainersVolume(
		repository2.NewStorePostgres(dbConnManager),
		repository2.NewStockPostgres(dbConnManager),
		logger,
	)

	usecase4 := usecase3.NewRefreshGroupContainers(
		repository2.NewContainerPostgres(dbConnManager),
		repository2.NewStockPostgres(dbConnManager),
		logger,
	)

	usecase5 := usecase3.NewRefreshLocations(
		usecase2,
		usecase4,
	)

	service2 := service.NewAccountSequence(
		repository.NewAccountSequencePostgres(dbConnManager),
	)

	useCase := usecase.NewCreateStockContainer(
		dbConnManager,
		storage,
		service2,
		storageStock,
		service.NewStore(
			storage2,
			storage,
		),
		usecase5,
		eventEmitter,
	)

	service1 := service.NewContainer(storage, eventEmitter)

	controller := httpv1.NewContainer(
		requestExtendParser,
		responseSender,
		service1,
		useCase,
	)

	return controller, nil
}
