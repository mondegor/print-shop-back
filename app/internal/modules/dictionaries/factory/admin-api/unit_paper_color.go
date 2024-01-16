package factory

import (
	module "print-shop-back/internal/modules/dictionaries"
	http_v1 "print-shop-back/internal/modules/dictionaries/controller/http_v1/admin-api"
	entity "print-shop-back/internal/modules/dictionaries/entity/admin-api"
	"print-shop-back/internal/modules/dictionaries/factory"
	repository "print-shop-back/internal/modules/dictionaries/infrastructure/repository/admin-api"
	usecase "print-shop-back/internal/modules/dictionaries/usecase/admin-api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrcore"
)

func newUnitPaperColor(
	c *[]mrcore.HttpController,
	opts *factory.Options,
	section mrcore.ClientSection,
) error {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(entity.PaperColor{})

	if err != nil {
		return err
	}

	storage := repository.NewPaperColorPostgres(
		opts.PostgresAdapter,
		mrsql.NewBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			mrpostgres.NewSqlBuilderOrderByWithDefaultSort(metaOrderBy.DefaultSort()),
			mrpostgres.NewSqlBuilderPager(module.PageSizeMax),
		),
	)
	service := usecase.NewPaperColor(storage, opts.EventBox, opts.ServiceHelper)
	*c = append(*c, http_v1.NewPaperColor(section, service, metaOrderBy))

	return nil
}
