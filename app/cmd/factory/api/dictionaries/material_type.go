package dictionaries

import (
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/api/availability/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/materialtype/api/availability"
)

// InitMaterialTypeAvailabilityAPI - создаёт объект usecase.MaterialType.
func InitMaterialTypeAvailabilityAPI(opts app.Options) *usecase.MaterialType {
	mrlog.Info(opts.Logger, "Create and init dictionaries laminate type availability API")

	return availability.NewMaterialTypeAPI(opts.PostgresConnManager, opts.UseCaseErrorWrapper, opts.Tracer)
}
