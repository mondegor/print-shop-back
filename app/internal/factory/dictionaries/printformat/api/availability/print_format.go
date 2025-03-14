package availability

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore/mrapp"

	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/api/availability/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/api/availability/usecase"
)

// NewPrintFormat - создаёт объект PrintFormat.
func NewPrintFormat(client mrstorage.DBConnManager) *usecase.PrintFormat {
	return usecase.NewPrintFormat(
		repository.NewPrintFormatPostgres(client),
		mrapp.NewUseCaseErrorWrapper(),
	)
}
