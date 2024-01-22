package factory

import (
	"print-shop-back/internal/modules"

	"github.com/mondegor/go-components/mrorderer"
)

func NewOrdererAPI(opts *modules.Options) *mrorderer.Component {
	opts.Logger.Info("Create and init orderer component")

	itemOrdererStorage := mrorderer.NewRepository(opts.PostgresAdapter)

	return mrorderer.NewComponent(itemOrdererStorage, opts.EventBox)
}
