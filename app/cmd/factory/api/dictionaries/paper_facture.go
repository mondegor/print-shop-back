package dictionaries

import (
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/api/availability/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/paperfacture/api/availability"
)

// InitPaperFactureAvailabilityAPI - создаёт объект usecase.PaperFacture.
func InitPaperFactureAvailabilityAPI(opts app.Options) *usecase.PaperFacture {
	mrlog.Info(opts.Logger, "Create and init dictionaries paper facture availability API")

	return availability.NewPaperFactureAPI(opts.PostgresConnManager, opts.Tracer)
}
