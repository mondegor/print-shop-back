package pub

import (
	"context"

	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/pub/repository"
	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/pub/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/catalog/laminate"
)

func createUnitLaminate(ctx context.Context, opts laminate.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitLaminate(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitLaminate(_ context.Context, opts laminate.Options) (*httpv1.Laminate, error) { //nolint:unparam
	storage := repository.NewLaminatePostgres(
		opts.DBConnManager,
	)
	useCase := usecase.NewLaminate(storage, opts.UseCaseHelper)
	controller := httpv1.NewLaminate(
		opts.RequestParsers.Parser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
