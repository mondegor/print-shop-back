package availability

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore/mrapp"

	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/api/availability/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/api/availability/usecase"
)

// NewMaterialType - создаёт объект MaterialType.
func NewMaterialType(client mrstorage.DBConnManager) *usecase.MaterialType {
	return usecase.NewMaterialType(
		repository.NewMaterialTypePostgres(client),
		mrapp.NewUseCaseErrorWrapper(),
	)
}
