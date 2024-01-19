package factory_api

import (
	repository_api "print-shop-back/internal/modules/controls/infrastructure/repository/api"
	usecase_api "print-shop-back/internal/modules/controls/usecase/api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-webcore/mrtool"
)

func NewElementTemplate(conn *mrpostgres.ConnAdapter, serviceHelper *mrtool.ServiceHelper) *usecase_api.ElementTemplate {
	return usecase_api.NewElementTemplate(
		repository_api.NewElementTemplatePostgres(conn),
		serviceHelper,
	)
}