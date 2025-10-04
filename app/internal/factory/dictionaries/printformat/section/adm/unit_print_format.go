package adm

import (
	"github.com/mondegor/go-storage/mrpostgres/builder"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/adm/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/adm/repository"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/adm/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/printformat"
)

func createUnitPrintFormat(opts printformat.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitPrintFormat(opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitPrintFormat(opts printformat.Options) (*httpv1.PrintFormat, error) {
	entityMeta, err := mrsql.ParseEntity(opts.Logger, entity.PrintFormat{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewPrintFormatPostgres(
		opts.DBConnManager,
		builder.NewSQL(
			builder.WithSQLOrderByDefaultSort(entityMeta.MetaOrderBy().DefaultSort()),
			builder.WithSQLLimitMaxSize(opts.PageSizeMax),
		),
	)
	useCase := usecase.NewPrintFormat(storage, opts.EventEmitter, opts.UsecaseErrorWrapper)
	controller := httpv1.NewPrintFormat(
		opts.RequestParsers.ExtendParser,
		opts.ResponseSender,
		useCase,
		entityMeta.MetaOrderBy(),
	)

	return controller, nil
}
