package factory_api

import (
	repository_api "print-shop-back/internal/modules/dictionaries/infrastructure/repository/api"
	usecase_api "print-shop-back/internal/modules/dictionaries/usecase/api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-webcore/mrtool"
)

func NewPaperFacture(conn *mrpostgres.ConnAdapter, serviceHelper *mrtool.ServiceHelper) *usecase_api.PaperFacture {
	return usecase_api.NewPaperFacture(
		repository_api.NewPaperFacturePostgres(conn),
		serviceHelper,
	)
}
