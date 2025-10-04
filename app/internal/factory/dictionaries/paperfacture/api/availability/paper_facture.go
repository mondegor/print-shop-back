package availability

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrtrace"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/api/availability/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/api/availability/usecase"
)

// NewPaperFacture - создаёт объект PaperFacture.
func NewPaperFacture(client mrstorage.DBConnManager, errorWrapper mrerr.UseCaseErrorWrapper, trace mrtrace.Tracer) *usecase.PaperFacture {
	return usecase.NewPaperFacture(
		repository.NewPaperFacturePostgres(client),
		errorWrapper,
		trace,
	)
}
