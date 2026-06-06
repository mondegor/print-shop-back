package availability

import (
	"github.com/mondegor/go-storage/mrstorage"

	"print-shop-back/internal/adapter/trace"
	"print-shop-back/internal/dictionaries/paperfacture/api/availability/repository"
	"print-shop-back/internal/dictionaries/paperfacture/api/availability/usecase"
)

// NewPaperFactureAPI - создаёт объект PaperFacture.
func NewPaperFactureAPI(
	dbConnManager mrstorage.DBConnManager,
	tracer trace.Tracer,
) *usecase.PaperFacture {
	return usecase.NewPaperFacture(
		repository.NewPaperFacturePostgres(dbConnManager),
		tracer,
	)
}
