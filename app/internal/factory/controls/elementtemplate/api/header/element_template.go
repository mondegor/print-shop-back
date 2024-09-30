package header

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/api/header/repository"
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/api/header/usecase"
)

// NewElementTemplate - создаёт объект ElementTemplate.
func NewElementTemplate(client mrstorage.DBConnManager, errorWrapper mrcore.UseCaseErrorWrapper) *usecase.ElementTemplate {
	return usecase.NewElementTemplate(
		repository.NewElementTemplatePostgres(client),
		errorWrapper,
	)
}
