package repository

import (
	"context"
	module "print-shop-back/internal/modules/catalog/paper"
	entity "print-shop-back/internal/modules/catalog/paper/entity/public-api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
)

type (
	PaperPostgres struct {
		client mrstorage.DBConn
	}
)

func NewPaperPostgres(
	client mrstorage.DBConn,
) *PaperPostgres {
	return &PaperPostgres{
		client: client,
	}
}

func (re *PaperPostgres) Fetch(ctx context.Context, params entity.PaperParams) ([]entity.Paper, error) {
	sql := `
        SELECT
            paper_id,
			paper_article,
			paper_caption,
			color_id,
			facture_id,
			paper_length,
			paper_width,
			paper_density,
			paper_thickness,
			paper_sides
        FROM
            ` + module.DBSchema + `.papers
        WHERE
            paper_status = $1
        ORDER BY
            paper_caption ASC, paper_id ASC;`

	cursor, err := re.client.Query(
		ctx,
		sql,
		mrenum.ItemStatusEnabled,
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]entity.Paper, 0)

	for cursor.Next() {
		var row entity.Paper

		err = cursor.Scan(
			&row.ID,
			&row.Article,
			&row.Caption,
			&row.ColorID,
			&row.FactureID,
			&row.Length,
			&row.Width,
			&row.Density,
			&row.Thickness,
			&row.Sides,
		)

		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, cursor.Err()
}
