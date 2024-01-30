package factory_api

import (
	repository_api "print-shop-back/internal/modules/dictionaries/infrastructure/repository/api"
	usecase_api "print-shop-back/internal/modules/dictionaries/usecase/api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-webcore/mrcore"
)

func NewPaperFacture(conn *mrpostgres.ConnAdapter, usecaseHelper *mrcore.UsecaseHelper) *usecase_api.PaperFacture {
	return usecase_api.NewPaperFacture(
		repository_api.NewPaperFacturePostgres(conn),
		usecaseHelper,
	)
}
