package calculations

import (
	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/factory/calculations/queryhistory"
)

// NewQueryHistoryModuleOptions - создаёт объект calculations.Options.
func NewQueryHistoryModuleOptions(opts app.Options) (queryhistory.Options, error) {
	return queryhistory.Options{
		Logger:              opts.Logger,
		EventEmitter:        opts.EventEmitter,
		UsecaseErrorWrapper: opts.UsecaseErrorWrapper,
		DBConnManager:       opts.PostgresConnManager,
		RequestParsers: queryhistory.RequestParsers{
			ExtendParser: opts.RequestParsers.ExtendParser,
		},
		ResponseSender: opts.ResponseSenders.Sender,
	}, nil
}
