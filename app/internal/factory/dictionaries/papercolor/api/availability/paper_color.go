package availability

import (
	"github.com/mondegor/go-core/mrstorage"

	"print-shop-back/internal/adapter/trace"
	"print-shop-back/internal/dictionaries/papercolor/api/availability/repository"
	"print-shop-back/internal/dictionaries/papercolor/api/availability/usecase"
)

// NewPaperColorAPI - создаёт объект PaperColor.
func NewPaperColorAPI(
	dbConnManager mrstorage.DBConnManager,
	tracer trace.Tracer,
) *usecase.PaperColor {
	return usecase.NewPaperColor(
		repository.NewPaperColorPostgres(dbConnManager),
		tracer,
	)
}
