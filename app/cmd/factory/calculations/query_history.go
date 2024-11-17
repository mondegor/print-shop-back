package calculations

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore/mrapp"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/factory/calculations/queryhistory"
)

// NewQueryHistoryModuleOptions - создаёт объект calculations.Options.
func NewQueryHistoryModuleOptions(_ context.Context, opts app.Options) (queryhistory.Options, error) {
	return queryhistory.Options{
		EventEmitter:        opts.EventEmitter,
		UseCaseErrorWrapper: mrapp.NewUseCaseErrorWrapper(),
		DBConnManager:       opts.PostgresConnManager,
		RequestParsers: queryhistory.RequestParsers{
			ExtendParser: opts.RequestParsers.ExtendParser,
		},
		ResponseSender: opts.ResponseSenders.Sender,
	}, nil
}
