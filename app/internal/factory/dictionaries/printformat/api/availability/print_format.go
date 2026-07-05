package availability

import (
	"github.com/mondegor/go-core/mrstorage"

	"print-shop-back/internal/adapter/trace"
	"print-shop-back/internal/dictionaries/printformat/api/availability/repository"
	"print-shop-back/internal/dictionaries/printformat/api/availability/usecase"
)

// NewPrintFormatAPI - создаёт объект PrintFormat.
func NewPrintFormatAPI(
	dbConnManager mrstorage.DBConnManager,
	tracer trace.Tracer,
) *usecase.PrintFormat {
	return usecase.NewPrintFormat(
		repository.NewPrintFormatPostgres(dbConnManager),
		tracer,
	)
}
