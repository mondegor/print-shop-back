package availability

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/api/availability/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/api/availability/usecase"
)

// NewPaperColor - создаёт объект PaperColor.
func NewPaperColor(client mrstorage.DBConnManager, errorWrapper mrcore.UseCaseErrorWrapper) *usecase.PaperColor {
	return usecase.NewPaperColor(
		repository.NewPaperColorPostgres(client),
		errorWrapper,
	)
}
