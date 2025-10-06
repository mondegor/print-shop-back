package dictionaries

import (
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/api/availability/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/paperfacture/api/availability"
)

// TODO: переделать!!!!

// NewPaperFactureAvailabilityAPI - создаёт объект usecase.PaperFacture.
func NewPaperFactureAvailabilityAPI(opts app.Options) (*usecase.PaperFacture, error) {
	mrlog.Info(opts.Logger, "Create and init dictionaries paper facture availability API")

	return availability.NewPaperFactureAPI(opts.PostgresConnManager, opts.UseCaseErrorWrapper, opts.Tracer), nil
}
