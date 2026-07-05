package availability

import (
	"github.com/mondegor/go-core/mrstorage"

	"print-shop-back/internal/adapter/trace"
	"print-shop-back/internal/dictionaries/materialtype/api/availability/repository"
	"print-shop-back/internal/dictionaries/materialtype/api/availability/usecase"
)

// NewMaterialTypeAPI - создаёт объект MaterialType.
func NewMaterialTypeAPI(
	dbConnManager mrstorage.DBConnManager,
	tracer trace.Tracer,
) *usecase.MaterialType {
	return usecase.NewMaterialType(
		repository.NewMaterialTypePostgres(dbConnManager),
		tracer,
	)
}
