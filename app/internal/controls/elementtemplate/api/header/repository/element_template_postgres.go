package repository

import (
	"context"

	"github.com/mondegor/go-storage/mrstorage"

	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/module"
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/section/adm/entity"
	"github.com/mondegor/print-shop-back/pkg/controls/api"
)

type (
	// ElementTemplatePostgres - comment struct.
	ElementTemplatePostgres struct {
		client mrstorage.DBConnManager
	}
)

// NewElementTemplatePostgres - создаёт объект ElementTemplatePostgres.
func NewElementTemplatePostgres(client mrstorage.DBConnManager) *ElementTemplatePostgres {
	return &ElementTemplatePostgres{
		client: client,
	}
}

// FetchOneHead - comment method.
func (re *ElementTemplatePostgres) FetchOneHead(ctx context.Context, rowID uint64) (entity.ElementTemplateHead, error) {
	sql := `
        SELECT
			tag_version,
            param_name,
            template_caption,
			element_detailing,
			template_status
        FROM
            ` + module.DBTableNameElementTemplates + `
        WHERE
            template_id = $1 AND deleted_at IS NULL
        LIMIT 1;`

	row := entity.ElementTemplateHead{
		ElementTemplateDTO: api.ElementTemplateDTO{ID: rowID},
	}

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&row.TagVersion,
		&row.ParamName,
		&row.Caption,
		&row.Detailing,
		&row.Status,
	)

	return row, err
}
