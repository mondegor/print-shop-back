package factory

import (
	"context"
	"print-shop-back/internal"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-webcore/mrlog"
)

func NewOrdererAPI(ctx context.Context, opts app.Options) *mrorderer.Component {
	mrlog.Ctx(ctx).Info().Msg("Create and init orderer component")
	itemOrdererStorage := mrorderer.NewRepository(opts.PostgresAdapter)

	return mrorderer.NewComponent(itemOrdererStorage, opts.EventEmitter)
}
