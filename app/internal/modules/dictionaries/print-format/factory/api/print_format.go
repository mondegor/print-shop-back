package factory_api

import (
	repository_api "print-shop-back/internal/modules/dictionaries/print-format/infrastructure/repository/api"
	usecase_api "print-shop-back/internal/modules/dictionaries/print-format/usecase/api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-webcore/mrcore"
)

func NewPrintFormat(conn *mrpostgres.ConnAdapter, usecaseHelper *mrcore.UsecaseHelper) *usecase_api.PrintFormat {
	return usecase_api.NewPrintFormat(
		repository_api.NewPrintFormatPostgres(conn),
		usecaseHelper,
	)
}
