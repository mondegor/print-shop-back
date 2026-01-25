package availability

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrtrace"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/api/availability/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/api/availability/usecase"
)

// NewPaperFactureAPI - создаёт объект PaperFacture.
func NewPaperFactureAPI(
	dbConnManager mrstorage.DBConnManager,
	trace mrtrace.Tracer,
) *usecase.PaperFacture {
	return usecase.NewPaperFacture(
		repository.NewPaperFacturePostgres(dbConnManager),
		trace,
	)
}
