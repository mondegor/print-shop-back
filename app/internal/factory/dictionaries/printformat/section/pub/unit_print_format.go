package pub

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/pub/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/pub/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/printformat"

	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitPrintFormat(ctx context.Context, opts printformat.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitPrintFormat(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitPrintFormat(_ context.Context, opts printformat.Options) (*httpv1.PrintFormat, error) { //nolint:unparam
	storage := repository.NewPrintFormatPostgres(
		opts.DBConnManager,
	)
	useCase := usecase.NewPrintFormat(storage, opts.UsecaseHelper)
	controller := httpv1.NewPrintFormat(
		opts.RequestParsers.Parser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
