package repository

import (
	"context"
	module "print-shop-back/internal/modules/catalog/box"
	entity "print-shop-back/internal/modules/catalog/box/entity/public-api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
)

type (
	BoxPostgres struct {
		client mrstorage.DBConn
	}
)

func NewBoxPostgres(
	client mrstorage.DBConn,
) *BoxPostgres {
	return &BoxPostgres{
		client: client,
	}
}

func (re *BoxPostgres) Fetch(ctx context.Context, params entity.BoxParams) ([]entity.Box, error) {
	sql := `
        SELECT
			box_id,
			box_article,
			box_caption,
			box_length,
			box_width,
			box_depth
        FROM
            ` + module.DBSchema + `.boxes
        WHERE
            box_status = $1
        ORDER BY
            box_caption ASC, box_id ASC;`

	cursor, err := re.client.Query(
		ctx,
		sql,
		mrenum.ItemStatusEnabled,
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]entity.Box, 0)

	for cursor.Next() {
		var row entity.Box

		err = cursor.Scan(
			&row.ID,
			&row.Article,
			&row.Caption,
			&row.Length,
			&row.Width,
			&row.Depth,
		)

		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, cursor.Err()
}
