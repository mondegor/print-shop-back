package pub

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/pub/repository"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/pub/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/controls/submitform"

	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitSubmitForm(ctx context.Context, opts submitform.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitSubmitForm(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitSubmitForm(_ context.Context, opts submitform.Options) (*httpv1.SubmitForm, error) { //nolint:unparam
	storage := repository.NewSubmitFormPostgres(
		opts.DBConnManager,
	)
	useCase := usecase.NewSubmitForm(storage, opts.UsecaseHelper)
	controller := httpv1.NewSubmitForm(
		opts.RequestParsers.ModuleParser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
