package pub

import (
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/pub/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/pub/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/paperfacture"
)

func createUnitPaperFacture(opts paperfacture.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitPaperFacture(opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitPaperFacture(opts paperfacture.Options) (*httpv1.PaperFacture, error) { //nolint:unparam
	storage := repository.NewPaperFacturePostgres(
		opts.DBConnManager,
	)
	useCase := usecase.NewPaperFacture(storage, opts.UsecaseErrorWrapper)
	controller := httpv1.NewPaperFacture(
		opts.RequestParsers.Parser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
