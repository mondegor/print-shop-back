package repository_api

import (
	"context"
	module "print-shop-back/internal/modules/controls/element-template"
	entity "print-shop-back/internal/modules/controls/element-template/entity/admin-api"
	"print-shop-back/pkg/modules/controls"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	ElementTemplatePostgres struct {
		client mrstorage.DBConn
	}
)

func NewElementTemplatePostgres(
	client mrstorage.DBConn,
) *ElementTemplatePostgres {
	return &ElementTemplatePostgres{
		client: client,
	}
}

func (re *ElementTemplatePostgres) FetchOneHead(ctx context.Context, rowID mrtype.KeyInt32) (entity.ElementTemplateHead, error) {
	sql := `
        SELECT
			tag_version,
            param_name,
            template_caption,
			element_detailing,
			template_status
        FROM
            ` + module.DBSchema + `.element_templates
        WHERE
            template_id = $1 AND template_status <> $2
        LIMIT 1;`

	row := entity.ElementTemplateHead{
		ElementTemplateHead: controls.ElementTemplateHead{ID: rowID},
	}

	err := re.client.QueryRow(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&row.TagVersion,
		&row.ParamName,
		&row.Caption,
		&row.Detailing,
		&row.Status,
	)

	return row, err
}
