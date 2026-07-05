package dictionaries

import (
	"print-shop-back/internal/adapter/log"
	"print-shop-back/internal/app"
	"print-shop-back/internal/dictionaries/paperfacture/api/availability/usecase"
	"print-shop-back/internal/factory/dictionaries/paperfacture/api/availability"
)

// InitPaperFactureAvailabilityAPI - создаёт объект usecase.PaperFacture.
func InitPaperFactureAvailabilityAPI(opts app.Options) *usecase.PaperFacture {
	log.Info(opts.Logger, "Create and init dictionaries paper facture availability API")

	return availability.NewPaperFactureAPI(opts.PostgresConnManager, opts.Tracer)
}
