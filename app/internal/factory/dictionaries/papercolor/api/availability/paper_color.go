package availability

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrtrace"

	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/api/availability/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/api/availability/usecase"
)

// NewPaperColorAPI - создаёт объект PaperColor.
func NewPaperColorAPI(
	dbConnManager mrstorage.DBConnManager,
	trace mrtrace.Tracer,
) *usecase.PaperColor {
	return usecase.NewPaperColor(
		repository.NewPaperColorPostgres(dbConnManager),
		trace,
	)
}
