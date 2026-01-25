package availability

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrtrace"

	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/api/availability/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/api/availability/usecase"
)

// NewPrintFormatAPI - создаёт объект PrintFormat.
func NewPrintFormatAPI(
	dbConnManager mrstorage.DBConnManager,
	trace mrtrace.Tracer,
) *usecase.PrintFormat {
	return usecase.NewPrintFormat(
		repository.NewPrintFormatPostgres(dbConnManager),
		trace,
	)
}
