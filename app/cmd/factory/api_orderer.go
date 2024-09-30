package factory

import (
	"context"

	"github.com/mondegor/go-components/factory/mrsort"
	"github.com/mondegor/go-components/mrsort/component/orderer"
	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/print-shop-back/internal/app"
)

// NewOrdererAPI - создаёт объект orderer.Component.
func NewOrdererAPI(ctx context.Context, opts app.Options) *orderer.Component {
	mrlog.Ctx(ctx).Info().Msg("Create and init orderer component")

	return mrsort.NewComponentOrderer(
		mrsort.ComponentOptions{
			DBClient:     opts.PostgresConnManager,
			EventEmitter: opts.EventEmitter,
		},
	)
}
