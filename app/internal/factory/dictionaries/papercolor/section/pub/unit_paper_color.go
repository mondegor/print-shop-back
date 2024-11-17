package pub

import (
	"context"

	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/pub/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/pub/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/papercolor"
)

func createUnitPaperColor(ctx context.Context, opts papercolor.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitPaperColor(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitPaperColor(_ context.Context, opts papercolor.Options) (*httpv1.PaperColor, error) { //nolint:unparam
	storage := repository.NewPaperColorPostgres(
		opts.DBConnManager,
	)
	useCase := usecase.NewPaperColor(storage, opts.UseCaseErrorWrapper)
	controller := httpv1.NewPaperColor(
		opts.RequestParsers.Parser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
