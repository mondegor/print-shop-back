package calculations

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore/mrapp"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/factory/calculations/algo"
)

// NewAlgoModuleOptions - создаёт объект calculations.Options.
func NewAlgoModuleOptions(_ context.Context, opts app.Options) (algo.Options, error) {
	return algo.Options{
		EventEmitter:        opts.EventEmitter,
		UseCaseErrorWrapper: mrapp.NewUseCaseErrorWrapper(),
		DBConnManager:       opts.PostgresConnManager,
		RequestParsers: algo.RequestParsers{
			Validator: opts.RequestParsers.Validator,
		},
		ResponseSender: opts.ResponseSenders.Sender,
	}, nil
}
