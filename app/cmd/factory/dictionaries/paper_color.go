package dictionaries

import (
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/api/availability/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/papercolor"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/papercolor/api/availability"
)

// NewPaperColorModuleOptions - создаёт объект papercolor.Options.
func NewPaperColorModuleOptions(opts app.Options) papercolor.Options {
	return papercolor.Options{
		Logger:              opts.Logger,
		EventEmitter:        opts.EventEmitter,
		UsecaseErrorWrapper: opts.UsecaseErrorWrapper,
		DBConnManager:       opts.PostgresConnManager,
		RequestParsers: papercolor.RequestParsers{
			Parser:       opts.RequestParsers.Parser,
			ExtendParser: opts.RequestParsers.ExtendParser,
		},
		ResponseSender: opts.ResponseSenders.Sender,

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}
}

// NewPaperColorAvailabilityAPI - создаёт объект usecase.PaperColor.
func NewPaperColorAvailabilityAPI(opts app.Options) (*usecase.PaperColor, error) {
	mrlog.Info(opts.Logger, "Create and init dictionaries paper color availability API")

	return availability.NewPaperColor(opts.PostgresConnManager, opts.UsecaseErrorWrapper, opts.Tracer), nil
}
