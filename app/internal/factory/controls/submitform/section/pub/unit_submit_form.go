package pub

import (
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/pub/repository"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/pub/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/controls/submitform"
)

func createUnitSubmitForm(opts submitform.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitSubmitForm(opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitSubmitForm(opts submitform.Options) (*httpv1.SubmitForm, error) { //nolint:unparam
	storage := repository.NewSubmitFormPostgres(
		opts.DBConnManager,
	)
	useCase := usecase.NewSubmitForm(storage, opts.UsecaseErrorWrapper)
	controller := httpv1.NewSubmitForm(
		opts.RequestParsers.ModuleParser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
