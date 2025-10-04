package adm

import (
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/factory/controls/submitform"
)

func initUnitSubmitFormVersionEnvironment(opts submitform.Options) (formVersionOptions, error) { //nolint:unparam
	storage := repository.NewFormVersionPostgres(
		opts.DBConnManager,
	)

	return formVersionOptions{
		storage: storage,
	}, nil
}
