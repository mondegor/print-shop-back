package header

import (
	"github.com/mondegor/go-sysmess/mrstorage"

	"print-shop-back/internal/adapter/trace"
	"print-shop-back/internal/controls/elementtemplate/api/header/repository"
	"print-shop-back/internal/controls/elementtemplate/api/header/usecase"
)

// NewElementTemplate - создаёт объект ElementTemplate.
func NewElementTemplate(client mrstorage.DBConnManager, tracer trace.Tracer) *usecase.ElementTemplate {
	return usecase.NewElementTemplate(
		repository.NewElementTemplatePostgres(client),
		tracer,
	)
}
