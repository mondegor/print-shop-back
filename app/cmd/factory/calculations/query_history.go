package calculations

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/calculations/queryhistory/module"
	"github.com/mondegor/print-shop-back/internal/factory/calculations/queryhistory"

	"github.com/mondegor/go-webcore/mrcore/mrinit"
)

// NewQueryHistoryModuleOptions - создаёт объект calculations.Options.
func NewQueryHistoryModuleOptions(_ context.Context, opts app.Options) (queryhistory.Options, error) {
	return queryhistory.Options{
		EventEmitter:  opts.EventEmitter,
		UsecaseHelper: opts.UsecaseErrorWrapper,
		DBConnManager: opts.PostgresConnManager,
		RequestParsers: queryhistory.RequestParsers{
			ExtendParser: opts.RequestParsers.ExtendParser,
		},
		ResponseSender: opts.ResponseSenders.Sender,
	}, nil
}

// RegisterQueryHistoryErrors - comment func.
func RegisterQueryHistoryErrors(em *mrinit.ErrorManager) {
	em.RegisterList(mrinit.WrapProtoList(module.Errors()))
}
