package pub

import (
	"context"

	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/pub/repository"
	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/pub/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/catalog/paper"
)

func createUnitPaper(ctx context.Context, opts paper.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitPaper(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitPaper(_ context.Context, opts paper.Options) (*httpv1.Paper, error) { //nolint:unparam
	storage := repository.NewPaperPostgres(
		opts.DBConnManager,
	)
	useCase := usecase.NewPaper(storage, opts.UseCaseErrorWrapper)
	controller := httpv1.NewPaper(
		opts.RequestParsers.Parser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
