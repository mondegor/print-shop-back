package usr

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	repository2 "github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/back/repository"
	usecase3 "github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/back/usecase"
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/repository"
	service2 "github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/service"
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/transport/httpv1"
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/usecase"
	"github.com/mondegor/print-shop-back/pkg/transport/validate"
)

func initStockController(
	logger mrlog.Logger,
	eventEmitter mrevent.Emitter,
	dbConnManager mrstorage.DBConnManager,
	requestExtendParser *validate.ExtendParser,
	responseSender mrserver.ResponseSender,
) (mrserver.HttpController, error) {
	// :TODO: временно, вынести на уровень выше
	storage3 := repository.NewContainerPostgres(
		dbConnManager,
	)

	stock := repository.NewStockPostgres(
		dbConnManager,
	)

	service := service2.NewStock(stock)

	// :TODO: временно, вынести на уровень выше
	storage2 := repository.NewStorePostgres(
		dbConnManager,
	)

	// :TODO: временно, вынести на уровень выше
	usecase1 := usecase.NewAddMoreStockContainer(
		storage3,
		stock,
		service2.NewStore(storage2, storage3),
		usecase3.NewRefreshStoreContainersVolume(
			repository2.NewStorePostgres(dbConnManager),
			repository2.NewStockPostgres(dbConnManager),
			logger,
		),
		eventEmitter,
	)

	// :TODO: временно, вынести на уровень выше
	usecase2 := usecase.NewMoveStockContainer(
		stock,
		service2.NewStore(storage2, storage3),
		usecase3.NewRefreshLocations(
			usecase3.NewRefreshStoreContainersVolume(
				repository2.NewStorePostgres(dbConnManager),
				repository2.NewStockPostgres(dbConnManager),
				logger,
			),
			usecase3.NewRefreshGroupContainers(
				repository2.NewContainerPostgres(dbConnManager),
				repository2.NewStockPostgres(dbConnManager),
				logger,
			),
		),
		eventEmitter,
	)

	// :TODO: временно, вынести на уровень выше
	usecase4 := usecase.NewTransferStockContainer(
		stock,
		usecase3.NewRefreshStoreContainersVolume(
			repository2.NewStorePostgres(dbConnManager),
			repository2.NewStockPostgres(dbConnManager),
			logger,
		),
		eventEmitter,
	)

	controller := httpv1.NewStock(
		requestExtendParser,
		responseSender,
		service,
		usecase2,
		usecase1,
		usecase4,
	)

	return controller, nil
}
