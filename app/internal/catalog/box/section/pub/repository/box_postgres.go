package repository

import (
	"context"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/print-shop-back/internal/catalog/box/module"
	"github.com/mondegor/print-shop-back/internal/catalog/box/section/pub/entity"
)

type (
	// BoxPostgres - comment struct.
	BoxPostgres struct {
		client mrstorage.DBConnManager
	}
)

// NewBoxPostgres - создаёт объект BoxPostgres.
func NewBoxPostgres(client mrstorage.DBConnManager) *BoxPostgres {
	return &BoxPostgres{
		client: client,
	}
}

// Fetch - comment method.
func (re *BoxPostgres) Fetch(ctx context.Context, _ entity.BoxParams) ([]entity.Box, error) {
	sql := `
        SELECT
			box_id,
			box_article,
			box_caption,
			box_length,
			box_width,
			box_height,
			box_thickness,
			box_weight
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameBoxes + `
        WHERE
            box_status = $1 AND deleted_at IS NULL
        ORDER BY
            box_caption ASC, box_id ASC;`

	cursor, err := re.client.Conn(ctx).Query(
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
			&row.Height,
			&row.Thickness,
			&row.Weight,
		)
		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, cursor.Err()
}
