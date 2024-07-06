package dictionaries

import (
	"context"

	"github.com/mondegor/print-shop-back/pkg/dictionaries/api"

	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/api/availability/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/printformat"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/printformat/api/availability"
)

// NewPrintFormatModuleOptions - создаёт объект printformat.Options.
func NewPrintFormatModuleOptions(_ context.Context, opts app.Options) (printformat.Options, error) {
	printFormatDictionary, err := opts.Translator.Dictionary("dictionaries/print-formats")
	if err != nil {
		return printformat.Options{}, err
	}

	return printformat.Options{
		EventEmitter:  opts.EventEmitter,
		UsecaseHelper: opts.UsecaseErrorWrapper,
		DBConnManager: opts.PostgresConnManager,
		RequestParsers: printformat.RequestParsers{
			Parser:       opts.RequestParsers.Parser,
			ExtendParser: opts.RequestParsers.ExtendParser,
		},
		ResponseSender: opts.ResponseSenders.Sender,

		UnitPrintFormat: printformat.UnitPrintFormatOptions{
			Dictionary: printFormatDictionary,
		},

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}

// NewPrintFormatAvailabilityAPI - создаёт объект usecase.PrintFormat.
func NewPrintFormatAvailabilityAPI(ctx context.Context, opts app.Options) (*usecase.PrintFormat, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init dictionaries print format availability API")

	return availability.NewPrintFormat(opts.PostgresConnManager, opts.UsecaseErrorWrapper), nil
}

// RegisterPrintFormatErrors - comment func.
func RegisterPrintFormatErrors(em *mrinit.ErrorManager) {
	em.RegisterList(mrinit.WrapProtoList(api.PrintFormatErrors()))
}
