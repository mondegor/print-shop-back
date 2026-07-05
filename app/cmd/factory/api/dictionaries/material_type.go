package dictionaries

import (
	"print-shop-back/internal/adapter/log"
	"print-shop-back/internal/app"
	"print-shop-back/internal/dictionaries/materialtype/api/availability/usecase"
	"print-shop-back/internal/factory/dictionaries/materialtype/api/availability"
)

// InitMaterialTypeAvailabilityAPI - создаёт объект usecase.MaterialType.
func InitMaterialTypeAvailabilityAPI(opts app.Options) *usecase.MaterialType {
	log.Info(opts.Logger, "Create and init dictionaries laminate type availability API")

	return availability.NewMaterialTypeAPI(opts.PostgresConnManager, opts.Tracer)
}
