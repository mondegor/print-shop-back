package factory

import (
	"context"
	module "print-shop-back/internal/modules/dictionaries/paper-facture"
	http_v1 "print-shop-back/internal/modules/dictionaries/paper-facture/controller/http_v1/admin-api"
	entity "print-shop-back/internal/modules/dictionaries/paper-facture/entity/admin-api"
	"print-shop-back/internal/modules/dictionaries/paper-facture/factory"
	repository "print-shop-back/internal/modules/dictionaries/paper-facture/infrastructure/repository/admin-api"
	usecase "print-shop-back/internal/modules/dictionaries/paper-facture/usecase/admin-api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitPaperFacture(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitPaperFacture(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitPaperFacture(ctx context.Context, opts factory.Options) (*http_v1.PaperFacture, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.PaperFacture{})

	if err != nil {
		return nil, err
	}

	storage := repository.NewPaperFacturePostgres(
		opts.PostgresAdapter,
		mrsql.NewBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			mrpostgres.NewSqlBuilderOrderByWithDefaultSort(ctx, metaOrderBy.DefaultSort()),
			mrpostgres.NewSqlBuilderPager(module.PageSizeMax),
		),
	)
	service := usecase.NewPaperFacture(storage, opts.EventEmitter, opts.UsecaseHelper)
	controller := http_v1.NewPaperFacture(
		opts.RequestParser,
		opts.ResponseSender,
		service,
		metaOrderBy,
	)

	return controller, nil
}
