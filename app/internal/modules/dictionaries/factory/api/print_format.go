package factory_api

import (
	repository_api "print-shop-back/internal/modules/dictionaries/infrastructure/repository/api"
	usecase_api "print-shop-back/internal/modules/dictionaries/usecase/api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-webcore/mrtool"
)

func NewPrintFormat(conn *mrpostgres.ConnAdapter, serviceHelper *mrtool.ServiceHelper) *usecase_api.PrintFormat {
	return usecase_api.NewPrintFormat(
		repository_api.NewPrintFormatPostgres(conn),
		serviceHelper,
	)
}
