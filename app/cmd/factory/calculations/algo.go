package calculations

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/module"
	"github.com/mondegor/print-shop-back/internal/factory/calculations/algo"

	"github.com/mondegor/go-webcore/mrcore/mrinit"

	"github.com/mondegor/print-shop-back/internal/app"
)

// NewAlgoModuleOptions - создаёт объект calculations.Options.
func NewAlgoModuleOptions(_ context.Context, opts app.Options) (algo.Options, error) {
	return algo.Options{
		EventEmitter:  opts.EventEmitter,
		UsecaseHelper: opts.UsecaseErrorWrapper,
		DBConnManager: opts.PostgresConnManager,
		RequestParsers: algo.RequestParsers{
			Validator: opts.RequestParsers.Validator,
		},
		ResponseSender: opts.ResponseSenders.Sender,
	}, nil
}

// RegisterAlgoErrors - comment func.
func RegisterAlgoErrors(em *mrinit.ErrorManager) {
	em.RegisterList(mrinit.WrapProtoList(module.Errors()))
}
