package dictionaries

import (
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/api/availability/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/paperfacture"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/paperfacture/api/availability"
)

// NewPaperFactureModuleOptions - создаёт объект paperfacture.Options.
func NewPaperFactureModuleOptions(opts app.Options) paperfacture.Options {
	return paperfacture.Options{
		Logger:              opts.Logger,
		EventEmitter:        opts.EventEmitter,
		UsecaseErrorWrapper: opts.UsecaseErrorWrapper,
		DBConnManager:       opts.PostgresConnManager,
		RequestParsers: paperfacture.RequestParsers{
			Parser:       opts.RequestParsers.Parser,
			ExtendParser: opts.RequestParsers.ExtendParser,
		},
		ResponseSender: opts.ResponseSenders.Sender,

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}
}

// NewPaperFactureAvailabilityAPI - создаёт объект usecase.PaperFacture.
func NewPaperFactureAvailabilityAPI(opts app.Options) (*usecase.PaperFacture, error) {
	mrlog.Info(opts.Logger, "Create and init dictionaries paper facture availability API")

	return availability.NewPaperFacture(opts.PostgresConnManager, opts.UsecaseErrorWrapper, opts.Tracer), nil
}
