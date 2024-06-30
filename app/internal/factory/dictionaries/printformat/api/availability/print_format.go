package availability

import (
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/api/availability/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/api/availability/usecase"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
)

// NewPrintFormat - создаёт объект PrintFormat.
func NewPrintFormat(client mrstorage.DBConnManager, errorWrapper mrcore.UsecaseErrorWrapper) *usecase.PrintFormat {
	return usecase.NewPrintFormat(
		repository.NewPrintFormatPostgres(client),
		errorWrapper,
	)
}
