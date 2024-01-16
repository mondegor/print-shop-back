package factory

import (
	module "print-shop-back/internal/modules/catalog"
	http_v1 "print-shop-back/internal/modules/catalog/controller/http_v1/admin-api"
	entity "print-shop-back/internal/modules/catalog/entity/admin-api"
	"print-shop-back/internal/modules/catalog/factory"
	repository "print-shop-back/internal/modules/catalog/infrastructure/repository/admin-api"
	usecase "print-shop-back/internal/modules/catalog/usecase/admin-api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrcore"
)

func newUnitPaper(
	c *[]mrcore.HttpController,
	opts *factory.Options,
	section mrcore.ClientSection,
) error {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(entity.Paper{})

	if err != nil {
		return err
	}

	entityMetaUpdate, err := mrsql.NewEntityMetaUpdate(entity.Paper{})

	if err != nil {
		return err
	}

	storage := repository.NewPaperPostgres(
		opts.PostgresAdapter,
		mrsql.NewBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			mrpostgres.NewSqlBuilderOrderByWithDefaultSort(metaOrderBy.DefaultSort()),
			mrpostgres.NewSqlBuilderPager(module.PageSizeMax),
		),
		mrsql.NewBuilderUpdateWithMeta(
			entityMetaUpdate,
			mrpostgres.NewSqlBuilderSet(),
		),
	)
	service := usecase.NewPaper(storage, opts.PaperColorAPI, opts.PaperFactureAPI, opts.EventBox, opts.ServiceHelper)
	*c = append(*c, http_v1.NewPaper(section, service, metaOrderBy))

	return nil
}
