package adm

import (
	"github.com/mondegor/go-storage/mrpostgres/builder"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/adm/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/adm/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/papercolor"
)

func createUnitPaperColor(opts papercolor.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitPaperColor(opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitPaperColor(opts papercolor.Options) (*httpv1.PaperColor, error) {
	entityMeta, err := mrsql.ParseEntity(opts.Logger, entity.PaperColor{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewPaperColorPostgres(
		opts.DBConnManager,
		builder.NewSQL(
			builder.WithSQLOrderByDefaultSort(entityMeta.MetaOrderBy().DefaultSort()),
			builder.WithSQLLimitMaxSize(opts.PageSizeMax),
		),
	)
	useCase := usecase.NewPaperColor(storage, opts.EventEmitter, opts.UsecaseErrorWrapper)
	controller := httpv1.NewPaperColor(
		opts.RequestParsers.ExtendParser,
		opts.ResponseSender,
		useCase,
		entityMeta.MetaOrderBy(),
	)

	return controller, nil
}
