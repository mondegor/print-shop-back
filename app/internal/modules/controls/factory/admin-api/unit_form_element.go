package factory

import (
	module "print-shop-back/internal/modules/controls"
	http_v1 "print-shop-back/internal/modules/controls/controller/http_v1/admin-api"
	entity "print-shop-back/internal/modules/controls/entity/admin-api"
	"print-shop-back/internal/modules/controls/factory"
	repository "print-shop-back/internal/modules/controls/infrastructure/repository/admin-api"
	usecase "print-shop-back/internal/modules/controls/usecase/admin-api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrcore"
)

func newUnitFormElement(
	c *[]mrcore.HttpController,
	opts *factory.Options,
	section mrcore.ClientSection,
) error {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(entity.FormElement{})

	if err != nil {
		return err
	}

	entityMetaUpdate, err := mrsql.NewEntityMetaUpdate(entity.FormElement{})

	if err != nil {
		return err
	}

	storage := repository.NewFormElementPostgres(
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
	service := usecase.NewFormElement(storage, opts.ElementTemplateAPI, opts.OrdererAPI, opts.EventBox, opts.ServiceHelper)
	*c = append(*c, http_v1.NewFormElement(section, service, metaOrderBy))

	return nil
}