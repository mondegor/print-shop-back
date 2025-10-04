package availability

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrtrace"

	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/api/availability/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/api/availability/usecase"
)

// NewPaperColor - создаёт объект PaperColor.
func NewPaperColor(client mrstorage.DBConnManager, errorWrapper mrerr.UseCaseErrorWrapper, trace mrtrace.Tracer) *usecase.PaperColor {
	return usecase.NewPaperColor(
		repository.NewPaperColorPostgres(client),
		errorWrapper,
		trace,
	)
}
