package header

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrtrace"

	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/api/header/repository"
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/api/header/usecase"
)

// NewElementTemplate - создаёт объект ElementTemplate.
func NewElementTemplate(client mrstorage.DBConnManager, trace mrtrace.Tracer) *usecase.ElementTemplate {
	return usecase.NewElementTemplate(
		repository.NewElementTemplatePostgres(client),
		trace,
	)
}
