package availability

import (
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/api/availability/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/api/availability/usecase"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
)

// NewPaperFacture - создаёт объект PaperFacture.
func NewPaperFacture(client mrstorage.DBConnManager, errorWrapper mrcore.UsecaseErrorWrapper) *usecase.PaperFacture {
	return usecase.NewPaperFacture(
		repository.NewPaperFacturePostgres(client),
		errorWrapper,
	)
}
