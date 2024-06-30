package header

import (
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/api/header/repository"
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/api/header/usecase"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
)

// NewElementTemplate - создаёт объект ElementTemplate.
func NewElementTemplate(client mrstorage.DBConnManager, errorWrapper mrcore.UsecaseErrorWrapper) *usecase.ElementTemplate {
	return usecase.NewElementTemplate(
		repository.NewElementTemplatePostgres(client),
		errorWrapper,
	)
}
