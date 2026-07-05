package dictionaries

import (
	"print-shop-back/internal/adapter/log"
	"print-shop-back/internal/app"
	"print-shop-back/internal/dictionaries/printformat/api/availability/usecase"
	"print-shop-back/internal/factory/dictionaries/printformat/api/availability"
)

// InitPrintFormatAvailabilityAPI - создаёт объект usecase.PrintFormat.
func InitPrintFormatAvailabilityAPI(opts app.Options) *usecase.PrintFormat {
	log.Info(opts.Logger, "Create and init dictionaries print format availability API")

	return availability.NewPrintFormatAPI(opts.PostgresConnManager, opts.Tracer)
}
