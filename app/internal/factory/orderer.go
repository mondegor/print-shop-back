package factory

import (
	"print-shop-back/config"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
)

func NewOrdererAPI(cfg *config.Config, client mrstorage.DBConn, logger mrcore.Logger, eventBox mrcore.EventBox) mrorderer.API {
	logger.Info("Create and init roles and permissions for modules access")

	itemOrdererStorage := mrorderer.NewRepository(client)

	return mrorderer.NewComponent(itemOrdererStorage, eventBox)
}
