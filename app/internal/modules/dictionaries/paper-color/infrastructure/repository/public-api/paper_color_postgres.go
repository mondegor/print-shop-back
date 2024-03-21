package repository

import (
	"context"
	module "print-shop-back/internal/modules/dictionaries/paper-color"
	entity "print-shop-back/internal/modules/dictionaries/paper-color/entity/public-api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
)

type (
	PaperColorPostgres struct {
		client mrstorage.DBConn
	}
)

func NewPaperColorPostgres(
	client mrstorage.DBConn,
) *PaperColorPostgres {
	return &PaperColorPostgres{
		client: client,
	}
}

func (re *PaperColorPostgres) Fetch(ctx context.Context, params entity.PaperColorParams) ([]entity.PaperColor, error) {
	sql := `
        SELECT
            color_id,
            color_caption
        FROM
            ` + module.DBSchema + `.paper_colors
        WHERE
            color_status = $1
        ORDER BY
            color_caption ASC, color_id ASC;`

	cursor, err := re.client.Query(
		ctx,
		sql,
		mrenum.ItemStatusEnabled,
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]entity.PaperColor, 0)

	for cursor.Next() {
		var row entity.PaperColor

		err = cursor.Scan(
			&row.ID,
			&row.Caption,
		)

		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, cursor.Err()
}
