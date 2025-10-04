package availability

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrtrace"

	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/api/availability/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/api/availability/usecase"
)

// NewMaterialType - создаёт объект MaterialType.
func NewMaterialType(client mrstorage.DBConnManager, errorWrapper mrerr.UseCaseErrorWrapper, trace mrtrace.Tracer) *usecase.MaterialType {
	return usecase.NewMaterialType(
		repository.NewMaterialTypePostgres(client),
		errorWrapper,
		trace,
	)
}
