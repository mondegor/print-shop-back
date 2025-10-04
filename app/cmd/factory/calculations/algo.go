package calculations

import (
	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/factory/calculations/algo"
)

// NewAlgoModuleOptions - создаёт объект calculations.Options.
func NewAlgoModuleOptions(opts app.Options) (algo.Options, error) {
	return algo.Options{
		Logger:        opts.Logger,
		EventEmitter:  opts.EventEmitter,
		DBConnManager: opts.PostgresConnManager,
		RequestParsers: algo.RequestParsers{
			Validator: opts.RequestParsers.Validator,
		},
		ResponseSender: opts.ResponseSenders.Sender,
	}, nil
}
