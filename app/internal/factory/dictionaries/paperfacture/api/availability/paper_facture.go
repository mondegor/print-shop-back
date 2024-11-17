package availability

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore/mrapp"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/api/availability/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/api/availability/usecase"
)

// NewPaperFacture - создаёт объект PaperFacture.
func NewPaperFacture(client mrstorage.DBConnManager) *usecase.PaperFacture {
	return usecase.NewPaperFacture(
		repository.NewPaperFacturePostgres(client),
		mrapp.NewUseCaseErrorWrapper(),
	)
}
