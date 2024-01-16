package repository_api

import (
	"context"
	module "print-shop-back/internal/modules/controls"
	entity "print-shop-back/internal/modules/controls/entity/api"

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

func (re *ElementTemplatePostgres) FetchHead(ctx context.Context, id mrtype.KeyInt32) (*entity.ElementTemplateHead, error) {
	sql := `
        SELECT
            param_name,
            template_caption
        FROM
            ` + module.UnitElementTemplateDBSchema + `.element_templates
        WHERE
            template_id = $1 AND template_status = $2
        LIMIT 1;`

	row := entity.ElementTemplateHead{ID: id}

	err := re.client.QueryRow(
		ctx,
		sql,
		id,
		mrenum.ItemStatusEnabled,
	).Scan(
		&row.ParamName,
		&row.Caption,
	)

	if err != nil {
		return nil, err
	}

	return &row, nil
}
