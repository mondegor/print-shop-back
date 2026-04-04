package usr

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/repository"
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/service"
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/transport/httpv1"
	"github.com/mondegor/print-shop-back/pkg/transport/validate"
)

func initStoreController(
	// logger mrlog.Logger,
	dbConnManager mrstorage.DBConnManager,
	requestExtendParser *validate.ExtendParser,
	responseSender mrserver.ResponseSender,
) (mrserver.HttpController, error) {
	storage := repository.NewStorePostgres(
		dbConnManager,
	)

	container := repository.NewContainerPostgres(
		dbConnManager,
	)

	serviceStore := service.NewStore(storage, container)

	controller := httpv1.NewStore(
		requestExtendParser,
		responseSender,
		serviceStore,
	)

	return controller, nil
}
