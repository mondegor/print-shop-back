package dictionaries

import (
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/api/availability/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/printformat"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/printformat/api/availability"
)

// NewPrintFormatModuleOptions - создаёт объект printformat.Options.
func NewPrintFormatModuleOptions(opts app.Options) printformat.Options {
	return printformat.Options{
		Logger:              opts.Logger,
		EventEmitter:        opts.EventEmitter,
		UsecaseErrorWrapper: opts.UsecaseErrorWrapper,
		DBConnManager:       opts.PostgresConnManager,
		RequestParsers: printformat.RequestParsers{
			Parser:       opts.RequestParsers.Parser,
			ExtendParser: opts.RequestParsers.ExtendParser,
		},
		ResponseSender: opts.ResponseSenders.Sender,

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}
}

// NewPrintFormatAvailabilityAPI - создаёт объект usecase.PrintFormat.
func NewPrintFormatAvailabilityAPI(opts app.Options) (*usecase.PrintFormat, error) {
	mrlog.Info(opts.Logger, "Create and init dictionaries print format availability API")

	return availability.NewPrintFormat(opts.PostgresConnManager, opts.UsecaseErrorWrapper, opts.Tracer), nil
}
