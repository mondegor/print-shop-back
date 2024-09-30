package pub

import (
	"context"

	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/pub/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/pub/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/materialtype"
)

func createUnitMaterialType(ctx context.Context, opts materialtype.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitMaterialType(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitMaterialType(_ context.Context, opts materialtype.Options) (*httpv1.MaterialType, error) { //nolint:unparam
	storage := repository.NewMaterialTypePostgres(
		opts.DBConnManager,
	)
	useCase := usecase.NewMaterialType(storage, opts.UseCaseHelper)
	controller := httpv1.NewMaterialType(
		opts.RequestParsers.Parser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
