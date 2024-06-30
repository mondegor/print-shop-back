package calculations

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/calculations/box/module"
	"github.com/mondegor/print-shop-back/internal/factory/calculations/box"

	"github.com/mondegor/go-webcore/mrcore/mrinit"
)

// NewBoxModuleOptions - создаёт объект calculations.Options.
func NewBoxModuleOptions(_ context.Context, opts app.Options) (box.Options, error) {
	return box.Options{
		EventEmitter:  opts.EventEmitter,
		UsecaseHelper: opts.UsecaseErrorWrapper,
		DBConnManager: opts.PostgresConnManager,
		RequestParsers: box.RequestParsers{
			ExtendParser: opts.RequestParsers.ExtendParser,
		},
		ResponseSender: opts.ResponseSenders.Sender,
	}, nil
}

// RegisterBoxErrors - comment func.
func RegisterBoxErrors(em *mrinit.ErrorManager) {
	em.RegisterList(mrinit.WrapProtoList(module.Errors()))
}
