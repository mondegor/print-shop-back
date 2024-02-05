package factory_api

import (
	repository_api "print-shop-back/internal/modules/dictionaries/paper-color/infrastructure/repository/api"
	usecase_api "print-shop-back/internal/modules/dictionaries/paper-color/usecase/api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-webcore/mrcore"
)

func NewPaperColor(conn *mrpostgres.ConnAdapter, usecaseHelper *mrcore.UsecaseHelper) *usecase_api.PaperColor {
	return usecase_api.NewPaperColor(
		repository_api.NewPaperColorPostgres(conn),
		usecaseHelper,
	)
}
