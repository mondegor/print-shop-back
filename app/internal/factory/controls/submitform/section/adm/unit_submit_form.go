package adm

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/controls/submitform"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
)

func initUnitSubmitFormEnvironment(ctx context.Context, opts submitform.Options) (submitFormOptions, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.SubmitForm{})
	if err != nil {
		return submitFormOptions{}, err
	}

	entityMetaUpdate, err := mrsql.NewEntityMetaUpdate(ctx, entity.SubmitForm{})
	if err != nil {
		return submitFormOptions{}, err
	}

	storage := repository.NewSubmitFormPostgres(
		opts.DBConnManager,
		mrpostgres.NewSQLBuilderSelect(
			mrpostgres.NewSQLBuilderWhere(),
			mrpostgres.NewSQLBuilderOrderBy(ctx, metaOrderBy.DefaultSort()),
			mrpostgres.NewSQLBuilderLimit(opts.PageSizeMax),
		),
		mrpostgres.NewSQLBuilderUpdateWithMeta(
			entityMetaUpdate,
			mrpostgres.NewSQLBuilderSet(),
			nil,
		),
	)

	return submitFormOptions{
		metaOrderBy: metaOrderBy,
		storage:     storage,
	}, nil
}

func createUnitSubmitForm(ctx context.Context, opts moduleOptions) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitSubmitForm(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitSubmitForm(_ context.Context, opts moduleOptions) (*httpv1.SubmitForm, error) { //nolint:unparam
	useCase := usecase.NewSubmitForm(
		opts.submitForm.storage,
		opts.formElement.storage,
		opts.formVersion.storage,
		opts.EventEmitter,
		opts.UsecaseHelper,
	)
	useCaseVersion := usecase.NewFormVersion(
		opts.formVersion.storage,
		useCase,
		usecase.NewFormCompilerJson(),
		opts.Locker,
		opts.EventEmitter,
		opts.UsecaseHelper,
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
