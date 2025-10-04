package adm

import (
	"github.com/mondegor/go-storage/mrpostgres/builder"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/controls/submitform"
)

func initUnitSubmitFormEnvironment(opts submitform.Options) (submitFormOptions, error) {
	entityMeta, err := mrsql.ParseEntity(opts.Logger, entity.SubmitForm{})
	if err != nil {
		return submitFormOptions{}, err
	}

	storage := repository.NewSubmitFormPostgres(
		opts.DBConnManager,
		builder.NewSQL(
			builder.WithSQLSetMetaEntity(entityMeta.MetaUpdate()),
			builder.WithSQLOrderByDefaultSort(entityMeta.MetaOrderBy().DefaultSort()),
			builder.WithSQLLimitMaxSize(opts.PageSizeMax),
		),
	)

	return submitFormOptions{
		metaOrderBy: entityMeta.MetaOrderBy(),
		storage:     storage,
	}, nil
}

func createUnitSubmitForm(opts moduleOptions) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitSubmitForm(opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitSubmitForm(opts moduleOptions) (*httpv1.SubmitForm, error) { //nolint:unparam
	useCase := usecase.NewSubmitForm(
		opts.submitForm.storage,
		opts.formElement.storage,
		opts.formVersion.storage,
		opts.EventEmitter,
		opts.UsecaseErrorWrapper,
	)
	useCaseVersion := usecase.NewFormVersion(
		opts.formVersion.storage,
		useCase,
		usecase.NewFormCompilerJson(),
		opts.Locker,
		opts.EventEmitter,
		opts.UsecaseErrorWrapper,
	)
	controller := httpv1.NewSubmitForm(
		opts.RequestParsers.ModuleParser,
		opts.ResponseSender,
		useCase,
		useCaseVersion,
		opts.submitForm.metaOrderBy,
	)

	return controller, nil
}
