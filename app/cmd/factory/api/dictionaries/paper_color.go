package dictionaries

import (
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/api/availability/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/papercolor/api/availability"
)

// InitPaperColorAvailabilityAPI - создаёт объект usecase.PaperColor.
func InitPaperColorAvailabilityAPI(opts app.Options) *usecase.PaperColor {
	mrlog.Info(opts.Logger, "Create and init dictionaries paper color availability API")

	return availability.NewPaperColorAPI(opts.PostgresConnManager, opts.UseCaseErrorWrapper, opts.Tracer)
}
