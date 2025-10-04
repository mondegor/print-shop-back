package pub

import (
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/pub/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/pub/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/materialtype"
)

func createUnitMaterialType(opts materialtype.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitMaterialType(opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitMaterialType(opts materialtype.Options) (*httpv1.MaterialType, error) { //nolint:unparam
	storage := repository.NewMaterialTypePostgres(
		opts.DBConnManager,
	)
	useCase := usecase.NewMaterialType(storage, opts.UsecaseErrorWrapper)
	controller := httpv1.NewMaterialType(
		opts.RequestParsers.Parser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
