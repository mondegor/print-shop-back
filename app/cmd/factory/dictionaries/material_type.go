package dictionaries

import (
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/api/availability/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/materialtype/api/availability"
)

// TODO: переделать!!!!

// NewMaterialTypeAvailabilityAPI - создаёт объект usecase.MaterialType.
func NewMaterialTypeAvailabilityAPI(opts app.Options) (*usecase.MaterialType, error) {
	mrlog.Info(opts.Logger, "Create and init dictionaries laminate type availability API")

	return availability.NewMaterialTypeAPI(opts.PostgresConnManager, opts.UseCaseErrorWrapper, opts.Tracer), nil
}
