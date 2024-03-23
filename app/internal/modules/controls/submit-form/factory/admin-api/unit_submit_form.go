package factory

import (
	"context"
	http_v1 "print-shop-back/internal/modules/controls/submit-form/controller/http_v1/admin-api"
	entity "print-shop-back/internal/modules/controls/submit-form/entity/admin-api"
	"print-shop-back/internal/modules/controls/submit-form/factory"
	repository "print-shop-back/internal/modules/controls/submit-form/infrastructure/repository/admin-api"
	usecase "print-shop-back/internal/modules/controls/submit-form/usecase/admin-api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresponse"
)

func initUnitSubmitFormEnvironment(ctx context.Context, opts factory.Options) (submitFormOptions, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.SubmitForm{})

	if err != nil {
		return submitFormOptions{}, err
	}

	entityMetaUpdate, err := mrsql.NewEntityMetaUpdate(ctx, entity.SubmitForm{})

	if err != nil {
		return submitFormOptions{}, err
	}

	storage := repository.NewSubmitFormPostgres(
		opts.PostgresAdapter,
		mrpostgres.NewSqlBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			mrpostgres.NewSqlBuilderOrderBy(ctx, metaOrderBy.DefaultSort()),
			mrpostgres.NewSqlBuilderPager(opts.PageSizeMax),
		),
		mrpostgres.NewSqlBuilderUpdateWithMeta(
			entityMetaUpdate,
			mrpostgres.NewSqlBuilderSet(),
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

func newUnitSubmitForm(ctx context.Context, opts moduleOptions) (*http_v1.SubmitForm, error) {
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
	controller := http_v1.NewSubmitForm(
		opts.RequestParser,
		mrresponse.NewFileSender(opts.ResponseSender),
		useCase,
		useCaseVersion,
		opts.submitForm.metaOrderBy,
	)

	return controller, nil
}
