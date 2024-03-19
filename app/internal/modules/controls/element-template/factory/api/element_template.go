package factory_api

import (
	repository_api "print-shop-back/internal/modules/controls/element-template/infrastructure/repository/api"
	usecase_api "print-shop-back/internal/modules/controls/element-template/usecase/api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-webcore/mrcore"
)

func NewElementTemplate(conn *mrpostgres.ConnAdapter, usecaseHelper *mrcore.UsecaseHelper) *usecase_api.ElementTemplate {
	return usecase_api.NewElementTemplate(
		repository_api.NewElementTemplatePostgres(conn),
		usecaseHelper,
	)
}
