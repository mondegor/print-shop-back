package catalog

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore/mrinit"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/catalog/box/module"
	"github.com/mondegor/print-shop-back/internal/factory/catalog/box"
)

// NewBoxModuleOptions - создаёт объект box.Options.
func NewBoxModuleOptions(_ context.Context, opts app.Options) (box.Options, error) {
	boxDictionary, err := opts.Translator.Dictionary("catalog/boxes")
	if err != nil {
		return box.Options{}, err
	}

	return box.Options{
		EventEmitter:  opts.EventEmitter,
		UseCaseHelper: opts.UseCaseErrorWrapper,
		DBConnManager: opts.PostgresConnManager,
		RequestParsers: box.RequestParsers{
			Parser:       opts.RequestParsers.Parser,
			ExtendParser: opts.RequestParsers.ExtendParser,
		},
		ResponseSender: opts.ResponseSenders.Sender,

		UnitBox: box.UnitBoxOptions{
			Dictionary: boxDictionary,
		},

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}

// RegisterBoxErrors - comment func.
func RegisterBoxErrors(em *mrinit.ErrorManager) {
	em.RegisterList(mrinit.WrapProtoList(module.Errors()))
}
