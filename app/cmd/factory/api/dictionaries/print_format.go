package dictionaries

import (
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/api/availability/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/printformat/api/availability"
)

// InitPrintFormatAvailabilityAPI - создаёт объект usecase.PrintFormat.
func InitPrintFormatAvailabilityAPI(opts app.Options) *usecase.PrintFormat {
	mrlog.Info(opts.Logger, "Create and init dictionaries print format availability API")

	return availability.NewPrintFormatAPI(opts.PostgresConnManager, opts.Tracer)
}
