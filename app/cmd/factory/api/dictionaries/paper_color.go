package dictionaries

import (
	"print-shop-back/internal/adapter/log"
	"print-shop-back/internal/app"
	"print-shop-back/internal/dictionaries/papercolor/api/availability/usecase"
	"print-shop-back/internal/factory/dictionaries/papercolor/api/availability"
)

// InitPaperColorAvailabilityAPI - создаёт объект usecase.PaperColor.
func InitPaperColorAvailabilityAPI(opts app.Options) *usecase.PaperColor {
	log.Info(opts.Logger, "Create and init dictionaries paper color availability API")

	return availability.NewPaperColorAPI(opts.PostgresConnManager, opts.Tracer)
}
