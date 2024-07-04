package pub

import (
	"context"

	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/queryhistory/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/calculations/queryhistory/section/pub/repository"
	"github.com/mondegor/print-shop-back/internal/calculations/queryhistory/section/pub/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/calculations/queryhistory"
)

func createUnitCalcResult(ctx context.Context, opts queryhistory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitCalcResult(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitCalcResult(_ context.Context, opts queryhistory.Options) (*httpv1.QueryHistory, error) { //nolint:unparam
	storage := repository.NewQueryHistoryPostgres(
		opts.DBConnManager,
	)
	useCase := usecase.NewQueryHistory(storage, opts.EventEmitter, opts.UsecaseHelper)
	controller := httpv1.NewQueryHistory(
		opts.RequestParsers.ExtendParser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
