package factory_api

import (
	repository_api "print-shop-back/internal/modules/dictionaries/laminate-type/infrastructure/repository/api"
	usecase_api "print-shop-back/internal/modules/dictionaries/laminate-type/usecase/api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-webcore/mrcore"
)

func NewLaminateType(conn *mrpostgres.ConnAdapter, usecaseHelper *mrcore.UsecaseHelper) *usecase_api.LaminateType {
	return usecase_api.NewLaminateType(
		repository_api.NewLaminateTypePostgres(conn),
		usecaseHelper,
	)
}
