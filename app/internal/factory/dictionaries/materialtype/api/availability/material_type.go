package availability

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/api/availability/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/api/availability/usecase"
)

// NewMaterialType - создаёт объект MaterialType.
func NewMaterialType(client mrstorage.DBConnManager, errorWrapper mrcore.UseCaseErrorWrapper) *usecase.MaterialType {
	return usecase.NewMaterialType(
		repository.NewMaterialTypePostgres(client),
		errorWrapper,
	)
}
