package availability

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrtrace"

	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/api/availability/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/api/availability/usecase"
)

// NewMaterialTypeAPI - создаёт объект MaterialType.
func NewMaterialTypeAPI(
	dbConnManager mrstorage.DBConnManager,
	trace mrtrace.Tracer,
) *usecase.MaterialType {
	return usecase.NewMaterialType(
		repository.NewMaterialTypePostgres(dbConnManager),
		trace,
	)
}
